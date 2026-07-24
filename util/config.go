package util

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

const (
	defaultAccessTokenDuration  = 15 * time.Minute
	defaultRefreshTokenDuration = 24 * time.Hour
)

type Config struct {
	Environment          string
	DBSource             string
	RedisAddress         string
	HTTPServerAddress    string
	GRPCServerAddress    string
	TokenSymmetricKey    string
	SmtpName             string
	SmtpUser             string
	SmtpPass             string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = ".env"
	}

	err := godotenv.Load(path)
	if err != nil {
		log.Warn().Err(err).Msg("error loading .env file, relying on system ENV")
	}

	environment, exist := os.LookupEnv("ENVIRONMENT")
	if !exist {
		return nil, errors.New("environments are not configured")
	}

	return &Config{
		Environment:          environment,
		DBSource:             getEnv("DB_SOURCE", ""),
		RedisAddress:         getEnv("REDIS_ADDRESS", ""),
		HTTPServerAddress:    getEnv("HTTP_SERVER_ADDRESS", "localhost:8080"),
		GRPCServerAddress:    getEnv("GRPC_SERVER_ADDRESS", "localhost:9090"),
		TokenSymmetricKey:    getEnv("TOKEN_SYMMETRIC_KEY", ""),
		SmtpName:             getEnv("SMTP_NAME", ""),
		SmtpUser:             getEnv("SMTP_USER", ""),
		SmtpPass:             getEnv("SMTP_PASS", ""),
		AccessTokenDuration:  getDurationEnv("ACCESS_TOKEN_DURATION", defaultAccessTokenDuration),
		RefreshTokenDuration: getDurationEnv("REFRESH_TOKEN_DURATION", defaultRefreshTokenDuration),
	}, nil
}

func getEnv(key, fallback string) string {
	value, exist := os.LookupEnv(key)
	if !exist {
		log.Warn().Msgf("env %s is missing", key)
		return fallback
	}

	return value
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	value, exist := os.LookupEnv(key)
	if !exist {
		log.Warn().Msgf("env %s is missing", key)
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Error().Err(err).Msgf("failed to parse env %s; value %s", key, value)
		return fallback
	}

	return duration
}
