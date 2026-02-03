package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration
type Config struct {
	Port            string
	PushoverToken   string
	PushoverUserKey string
	PushoverRetry   time.Duration
	PushoverExpire  time.Duration
	PushoverPriority int
	LogLevel        string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	token := os.Getenv("PUSHOVER_APP_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("PUSHOVER_APP_TOKEN is required")
	}

	userKey := os.Getenv("PUSHOVER_USER_KEY")
	if userKey == "" {
		return nil, fmt.Errorf("PUSHOVER_USER_KEY is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	priority := 2 // Emergency priority by default
	if p := os.Getenv("PUSHOVER_PRIORITY"); p != "" {
		var err error
		priority, err = strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid PUSHOVER_PRIORITY: %w", err)
		}
	}

	retry := 60 * time.Second
	if r := os.Getenv("PUSHOVER_RETRY"); r != "" {
		seconds, err := strconv.Atoi(r)
		if err != nil {
			return nil, fmt.Errorf("invalid PUSHOVER_RETRY: %w", err)
		}
		retry = time.Duration(seconds) * time.Second
	}

	expire := 3600 * time.Second
	if e := os.Getenv("PUSHOVER_EXPIRE"); e != "" {
		seconds, err := strconv.Atoi(e)
		if err != nil {
			return nil, fmt.Errorf("invalid PUSHOVER_EXPIRE: %w", err)
		}
		expire = time.Duration(seconds) * time.Second
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	return &Config{
		Port:             port,
		PushoverToken:    token,
		PushoverUserKey:  userKey,
		PushoverRetry:    retry,
		PushoverExpire:   expire,
		PushoverPriority: priority,
		LogLevel:         logLevel,
	}, nil
}
