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
	"github.com/ncfex/dcart-auth/internal/adapters/secondary/memory"
	"github.com/ncfex/dcart-auth/internal/adapters/secondary/postgres"
	"github.com/ncfex/dcart-auth/internal/application/services"

	"github.com/ncfex/dcart-auth/internal/config"

	"github.com/ncfex/dcart-auth/pkg/httputil/response"
	"github.com/ncfex/dcart-auth/pkg/services/auth/tokens/jwt"
	"github.com/ncfex/dcart-auth/pkg/services/auth/tokens/refresh"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	postgresURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	db, err := postgres.NewDatabase(postgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// repo
	// userRepo := postgres.NewUserRepository(db)
	tokenRepo := postgres.NewTokenRepository(db, 24*7*time.Hour)
	eventStore := memory.NewInMemoryEventStore()

	// cqrs
	userCommandHandler := services.NewUserCommandHandler(eventStore)
	userQueryHandler := services.NewUserQueryHandler(eventStore)

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
		eventStore,
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
