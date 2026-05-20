package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	GinMode                string
	MongoURI               string
	MongoDB                string
	JWTSecret              string
	JWTExpiry              time.Duration
	RefreshExpiry          time.Duration
	BlockchainRPC          string
	ContractAddress        string
	DeployerPrivateKey     string
	PinataAPIKey           string
	PinataSecret           string
	PinataJWT              string
	SMTPHost               string
	SMTPPort               string
	SMTPUser               string
	SMTPPass               string
	SMTPFrom               string
	CORSOrigins            []string
	MaxUploadMB            int64
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	refreshDays, _ := strconv.Atoi(getEnv("REFRESH_TOKEN_EXPIRY_DAYS", "7"))
	maxMB, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_MB", "10"), 10, 64)

	return &Config{
		Port:               getEnv("PORT", "8080"),
		GinMode:            getEnv("GIN_MODE", "debug"),
		MongoURI:           getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDB:            getEnv("MONGODB_DB", "fir_management"),
		JWTSecret:          getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		JWTExpiry:          time.Duration(jwtHours) * time.Hour,
		RefreshExpiry:      time.Duration(refreshDays) * 24 * time.Hour,
		BlockchainRPC:      getEnv("BLOCKCHAIN_RPC_URL", "http://127.0.0.1:8545"),
		ContractAddress:    getEnv("CONTRACT_ADDRESS", ""),
		DeployerPrivateKey: getEnv("DEPLOYER_PRIVATE_KEY", ""),
		PinataAPIKey:       getEnv("PINATA_API_KEY", ""),
		PinataSecret:       getEnv("PINATA_SECRET_KEY", ""),
		PinataJWT:          getEnv("PINATA_JWT", ""),
		SMTPHost:           getEnv("SMTP_HOST", ""),
		SMTPPort:           getEnv("SMTP_PORT", "587"),
		SMTPUser:           getEnv("SMTP_USER", ""),
		SMTPPass:           getEnv("SMTP_PASS", ""),
		SMTPFrom:           getEnv("SMTP_FROM", "noreply@fir-management.local"),
		CORSOrigins:        strings.Split(getEnv("CORS_ORIGINS", "*"), ","),
		MaxUploadMB:        maxMB,
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
