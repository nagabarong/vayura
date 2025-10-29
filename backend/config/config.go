package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Global database instance (legacy support)
var DB *gorm.DB

// Upload image
var UploadDir string

// Config holds application configuration
type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
	Storage  StorageConfig
}

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

type ServerConfig struct {
	Port string
}

type StorageConfig struct {
	UploadDir string
}

// Load reads configuration from environment variables
func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found, using system environment")
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", ""),
			DBName:   getEnvOrDefault("DB_NAME", "vayura"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
		},
		JWT: JWTConfig{
			Secret:    getEnvOrDefault("JWT_SECRET", "your-secret-key-change-in-production"),
			ExpiresIn: getDurationEnvOrDefault("JWT_EXPIRES_IN", "72h"),
		},
		Server: ServerConfig{
			Port: getEnvOrDefault("APP_PORT", "8080"),
		},
		Storage: StorageConfig{
			UploadDir: getEnvOrDefault("UPLOAD_DIR", "Uploads/avatars"),
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getDurationEnvOrDefault(key, defaultValue string) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		d, _ := time.ParseDuration(defaultValue)
		return d
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		// fallback to default on parse error
		d, _ = time.ParseDuration(defaultValue)
		return d
	}
	return d
}

// InitDB initializes database connection with legacy global variable
func InitDB() {
	cfg := Load()
	db, err := initDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	DB = db
	log.Println("✅ Database connected successfully")
}

// InitDatabase creates a new database connection
func InitDatabase(cfg DatabaseConfig) (*gorm.DB, error) {
	return initDatabase(cfg)
}

func initDatabase(cfg DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Set to logger.Info for SQL logging
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
