package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ncfex/dcart-auth/internal/adapters/primary/http/handlers"
	"github.com/ncfex/dcart-auth/internal/adapters/secondary/id"
	"github.com/ncfex/dcart-auth/internal/adapters/secondary/messaging/mqtest"
	"github.com/ncfex/dcart-auth/internal/adapters/secondary/persistence/mongodb"
	"github.com/ncfex/dcart-auth/internal/adapters/secondary/persistence/postgres"

	"github.com/ncfex/dcart-auth/internal/application/command"
	"github.com/ncfex/dcart-auth/internal/application/services"

	"github.com/ncfex/dcart-auth/internal/config"

	"github.com/ncfex/dcart-auth/pkg/httputil/response"
	"github.com/ncfex/dcart-auth/pkg/services/auth/tokens/jwt"
	"github.com/ncfex/dcart-auth/pkg/services/auth/tokens/refresh"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// write db
	// todo improve
	postgresURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	postgresDB, err := postgres.NewDatabase(postgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer postgresDB.Close()

	// read db
	mongoConfig := mongodb.Config{
		URI:            cfg.MongoURI,
		Database:       cfg.MongoCollection,
		ConnectTimeout: 10 * time.Second,
		MaxPoolSize:    100,
		MinPoolSize:    10,
	}

	mongoClient, err := mongodb.NewClient(mongoConfig)
	if err != nil {
		log.Fatal(err)
	}

	if err := mongoClient.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// persist
	tokenRepo := postgres.NewTokenRepository(postgresDB, 24*7*time.Hour)
	postgresEventStore := postgres.NewPostgresEventStore(postgresDB.DB)

	// test
	testPublisher := mqtest.NewEventPublisher()

	// id
	deterministicIDGen := id.NewDeterministicIDGenerator("dcart")

	// cqrs
	userCommandHandler := command.NewUserCommandHandler(
		postgresEventStore,
		testPublisher,
		deterministicIDGen,
	)

	// todo improve
	userQueryHandler := mongodb.NewUserQueryHandler(mongoClient.Database())

	// security
	jwtManager := jwt.NewJWTService("dcart", cfg.JwtSecret, time.Minute*15)
	refreshTokenGenerator := refresh.NewHexRefreshGenerator("dc_", 32)

	// app
	tokenSvc := services.NewTokenService(jwtManager, refreshTokenGenerator, tokenRepo)
	authService := services.NewAuthService(
		userCommandHandler,
		userQueryHandler,
		tokenSvc,
	)

	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	responder := response.NewHTTPResponder(logger)

	handler := handlers.NewHandler(
		logger,
		responder,
		authService,
		jwtManager,
		tokenRepo,
		postgresEventStore,
	)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler.RegisterRoutes(),

		// timeout
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		log.Printf("starting auth service on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not start server: %v", err)
		}
	}()

	<-stopSignal
	log.Println("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("http server shutting down: %v", err)
	}

	if err := postgresDB.Close(); err != nil {
		log.Printf("Error during database close: %v", err)
	}

	if err := mongoClient.Disconnect(ctx); err != nil {
		log.Printf("Error during database close: %v", err)
	}

	waitCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:
		log.Println("server stopped gracefully")
	case <-ctx.Done():
		log.Println("shutdown timed out")
	}
}
