package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Version          string // Version of the application
	ServiceName      string // Name of the service
	ServiceHost      string // Host of the service
	ServicePort      string // Port of the service
	ServiceHTTPS     string // HTTPS of the service
	TelegramBotToken string // Telegram bot token

	RedisHost   string
	RedisPrefix string

	LoggingLevel        int    // Logging level (integer value)
	LoggingChannel      string // Logging channel (file, database, etc.)
	LoggingEndpoint     string // Logging endpoint (file path, database URL, etc.)
	PriceTrackerService string
}

func NewConfig(envPath string) *Config {
	c := Config{}
	if envPath == "" {
		envPath = ".env"
	}
	c.initialize(envPath)
	return &c
}

func (c *Config) initialize(envPath string) {
	// Load config
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Environment file missed. Err: %s", err)
	}

	c.Version = os.Getenv("VERSION")
	if c.Version == "" {
		log.Panicln("VERSION not specified")
	}

	c.ServiceName = os.Getenv("SERVICE_NAME")
	if c.ServiceName == "" {
		log.Panicln("SERVICE_NAME not specified")
	}

	c.ServiceHost = os.Getenv("SERVICE_HOST")
	c.ServicePort = os.Getenv("SERVICE_PORT")
	c.ServiceHTTPS = os.Getenv("SERVICE_HTTPS")

	c.TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	if c.TelegramBotToken == "" {
		log.Panicln("TELEGRAM_BOT_TOKEN not specified")
	}

	c.RedisHost = os.Getenv("REDIS_HOST")
	c.RedisPrefix = os.Getenv("REDIS_KEY_PREFIX")

	ll, err := strconv.Atoi(os.Getenv("LOGGING_LEVEL"))
	if err != nil {
		ll = 3 // Default value
	}
	c.LoggingLevel = ll
	c.LoggingEndpoint = os.Getenv("LOGGING_ENDPOINT")
	c.LoggingChannel = os.Getenv("LOGGING_CHANNEL")

	c.PriceTrackerService = os.Getenv("PRICE_TRACKER_SERVICE")
	if c.ServiceName == "" {
		log.Panicln("PRICE_TRACKER_SERVICE not specified")
	}
}

func (c *Config) GetServiceEndpoint(serviceName string) (string, error) {
	switch serviceName {
	case c.PriceTrackerService:
		return c.PriceTrackerService, nil
	default:
		return "", fmt.Errorf("service endpoint not found for service: %s", serviceName)
	}
}
