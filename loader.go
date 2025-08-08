package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// LoadStrategy defines how configuration should be loaded
type LoadStrategy int

const (
	EnvironmentStrategy LoadStrategy = iota
	FileStrategy
	HybridStrategy
)

// Loader provides methods to load configuration
type Loader struct {
	viper *viper.Viper
}

// NewLoader creates a new configuration loader
func NewLoader() *Loader {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	return &Loader{
		viper: v,
	}
}

// LoadFromFile loads configuration from a file
func (l *Loader) LoadFromFile(configPath string) (*Config, error) {
	l.viper.SetConfigFile(configPath)

	if err := l.viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := l.viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// LoadFromEnvironment loads configuration from environment variables
func (l *Loader) LoadFromEnvironment() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "app"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			MaxConns: getIntEnv("DB_MAX_CONNS", 10),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		Log: LogConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Format:     getEnv("LOG_FORMAT", "json"),
			OutputPath: getEnv("LOG_OUTPUT_PATH", ""),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key"),
			Expiration: getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
			Issuer:     getEnv("JWT_ISSUER", "app"),
		},
		Email: EmailConfig{
			Host:     getEnv("EMAIL_HOST", ""),
			Port:     getIntEnv("EMAIL_PORT", 587),
			Username: getEnv("EMAIL_USERNAME", ""),
			Password: getEnv("EMAIL_PASSWORD", ""),
			From:     getEnv("EMAIL_FROM", ""),
		},
		App: AppConfig{
			Name:        getEnv("APP_NAME", "app"),
			Environment: getEnv("APP_ENVIRONMENT", "development"),
			Version:     getEnv("APP_VERSION", "1.0.0"),
			Debug:       getBoolEnv("APP_DEBUG", false),
		},
	}

	return config, nil
}

// Load loads configuration using the specified strategy
func (l *Loader) Load(strategy LoadStrategy) (*Config, error) {
	switch strategy {
	case FileStrategy:
		configPath := getEnv("CONFIG_PATH", "config.yaml")
		return l.LoadFromFile(configPath)
	case EnvironmentStrategy:
		return l.LoadFromEnvironment()
	case HybridStrategy:
		// Try file first, fallback to environment
		if configPath := getEnv("CONFIG_PATH", ""); configPath != "" {
			if config, err := l.LoadFromFile(configPath); err == nil {
				return config, nil
			}
		}
		return l.LoadFromEnvironment()
	default:
		return l.LoadFromEnvironment()
	}
}

// Helper functions for environment variable handling
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := parseInt(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := parseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// Parse functions
func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

func parseBool(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", s)
	}
}
