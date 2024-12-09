package config

import (
	"os"

	"github.com/joho/godotenv"
)

// todo use viper
type config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string
	MongoURI         string
	MongoDatabase    string
	RedisHost        string
	RedisPort        string
	JwtSecret        string
	Port             string
}

func LoadConfig() (*config, error) {
	err := godotenv.Load()
	if err != nil {
		return &config{}, err
	}
	return &config{
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresDB:       getEnv("POSTGRES_DB", "authdb"),
		PostgresUser:     getEnv("POSTGRES_USER", "guest"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "Z3Vlc3Q="), // base64 guest
		MongoURI:         getEnv("MONGO_URI", "mongodb://dev_user:dev_password@localhost:27017"),
		MongoDatabase:    getEnv("MONGO_DATABASE", "authdb"),
		RedisHost:        getEnv("REDIS_HOST", "localhost"),
		RedisPort:        getEnv("REDIS_PORT", "6379"),
		Port:             getEnv("AUTH_SERVICE_PORT", "8080"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
