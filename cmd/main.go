package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ncfex/dcart/auth-service/internal/adapters/primary/http/handlers"
	"github.com/ncfex/dcart/auth-service/internal/adapters/primary/http/response"
	"github.com/ncfex/dcart/auth-service/internal/adapters/secondary/postgres"
	"github.com/ncfex/dcart/auth-service/internal/core/services/authentication"
	"github.com/ncfex/dcart/auth-service/internal/core/services/password"
	"github.com/ncfex/dcart/auth-service/internal/core/services/refresh"
	"github.com/ncfex/dcart/auth-service/internal/core/services/token"
	"github.com/ncfex/dcart/auth-service/internal/infrastructure/config"
	postgresDB "github.com/ncfex/dcart/auth-service/internal/infrastructure/database/postgres"
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

	db, err := postgresDB.NewDatabase(postgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// repo
	userRepo := postgres.NewUserRepository(db)
	tokenRepo := postgres.NewTokenRepository(db, 24*7*time.Hour)
	if err != nil {
		log.Fatalf("Failed to initialize token repository: %v", err)
	}

	passwordService := password.NewPasswordService(0)
	jwtService := token.NewJWTService("dcart", cfg.JwtSecret)
	hexTokenService := refresh.NewHexTokenService("dc_", 32)
	authService := authentication.NewAuthService(
		userRepo,
		tokenRepo,
		passwordService,
		jwtService,
		hexTokenService,
	)

	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	responder := response.NewHTTPResponder(logger)

	handler := handlers.NewHandler(
		logger,
		responder,
		authService,
		jwtService,
		tokenRepo,
		userRepo,
	)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler.Router(),
	}

	log.Printf("starting auth service on port %s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
