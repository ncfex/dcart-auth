package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ncfex/dcart-auth/internal/adapters/primary/http/handlers"
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
	userRepo := postgres.NewUserRepository(db)
	tokenRepo := postgres.NewTokenRepository(db, 24*7*time.Hour)

	jwtManager := jwt.NewJWTService("dcart", cfg.JwtSecret, time.Minute*15)
	refreshTokenGenerator := refresh.NewHexRefreshGenerator("dc_", 32)

	// app
	userSvc := services.NewUserService(userRepo)
	tokenSvc := services.NewTokenService(jwtManager, refreshTokenGenerator, tokenRepo)
	authService := services.NewAuthService(
		userSvc,
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
		userRepo,
	)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler.RegisterRoutes(),
	}

	log.Printf("starting auth service on port %s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
