package config

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Validator provides configuration validation functionality
type Validator struct {
	errors []string
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		errors: make([]string, 0),
	}
}

// Validate validates the entire configuration
func (v *Validator) Validate(config *Config) error {
	v.errors = make([]string, 0)

	v.validateServer(config.Server)
	v.validateDatabase(config.Database)
	v.validateRedis(config.Redis)
	v.validateLog(config.Log)
	v.validateJWT(config.JWT)
	v.validateEmail(config.Email)
	v.validateApp(config.App)

	if len(v.errors) > 0 {
		return &ValidationError{
			Errors: v.errors,
		}
	}

	return nil
}

// ValidationError represents validation errors
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("configuration validation failed: %s", strings.Join(e.Errors, "; "))
}

// validateServer validates server configuration
func (v *Validator) validateServer(config ServerConfig) {
	if config.Port == "" {
		v.errors = append(v.errors, "server port is required")
	} else {
		if _, err := strconv.Atoi(config.Port); err != nil {
			v.errors = append(v.errors, "server port must be a valid integer")
		}
	}

	if config.Host == "" {
		v.errors = append(v.errors, "server host is required")
	}

	if config.ReadTimeout <= 0 {
		v.errors = append(v.errors, "server read timeout must be positive")
	}

	if config.WriteTimeout <= 0 {
		v.errors = append(v.errors, "server write timeout must be positive")
	}

	if config.IdleTimeout <= 0 {
		v.errors = append(v.errors, "server idle timeout must be positive")
	}
}

// validateDatabase validates database configuration
func (v *Validator) validateDatabase(config DatabaseConfig) {
	if config.Host == "" {
		v.errors = append(v.errors, "database host is required")
	}

	if config.Port == "" {
		v.errors = append(v.errors, "database port is required")
	} else {
		if _, err := strconv.Atoi(config.Port); err != nil {
			v.errors = append(v.errors, "database port must be a valid integer")
		}
	}

	if config.User == "" {
		v.errors = append(v.errors, "database user is required")
	}

	if config.DBName == "" {
		v.errors = append(v.errors, "database name is required")
	}

	if config.MaxConns <= 0 {
		v.errors = append(v.errors, "database max connections must be positive")
	}

	// Validate SSL mode
	validSSLModes := []string{"disable", "require", "verify-ca", "verify-full"}
	valid := false
	for _, mode := range validSSLModes {
		if config.SSLMode == mode {
			valid = true
			break
		}
	}
	if !valid {
		v.errors = append(v.errors, fmt.Sprintf("database SSL mode must be one of: %s", strings.Join(validSSLModes, ", ")))
	}
}

// validateRedis validates Redis configuration
func (v *Validator) validateRedis(config RedisConfig) {
	if config.Host == "" {
		v.errors = append(v.errors, "redis host is required")
	}

	if config.Port == "" {
		v.errors = append(v.errors, "redis port is required")
	} else {
		if _, err := strconv.Atoi(config.Port); err != nil {
			v.errors = append(v.errors, "redis port must be a valid integer")
		}
	}

	if config.DB < 0 || config.DB > 15 {
		v.errors = append(v.errors, "redis database number must be between 0 and 15")
	}
}

// validateLog validates logging configuration
func (v *Validator) validateLog(config LogConfig) {
	validLevels := []string{"debug", "info", "warn", "warning", "error", "fatal", "panic"}
	valid := false
	for _, level := range validLevels {
		if strings.ToLower(config.Level) == level {
			valid = true
			break
		}
	}
	if !valid {
		v.errors = append(v.errors, fmt.Sprintf("log level must be one of: %s", strings.Join(validLevels, ", ")))
	}

	validFormats := []string{"json", "text", "console"}
	valid = false
	for _, format := range validFormats {
		if strings.ToLower(config.Format) == format {
			valid = true
			break
		}
	}
	if !valid {
		v.errors = append(v.errors, fmt.Sprintf("log format must be one of: %s", strings.Join(validFormats, ", ")))
	}
}

// validateJWT validates JWT configuration
func (v *Validator) validateJWT(config JWTConfig) {
	if config.Secret == "" {
		v.errors = append(v.errors, "JWT secret is required")
	} else if len(config.Secret) < 32 {
		v.errors = append(v.errors, "JWT secret must be at least 32 characters long")
	}

	if config.Expiration <= 0 {
		v.errors = append(v.errors, "JWT expiration must be positive")
	}

	if config.Issuer == "" {
		v.errors = append(v.errors, "JWT issuer is required")
	}
}

// validateEmail validates email configuration
func (v *Validator) validateEmail(config EmailConfig) {
	if config.Host != "" {
		if config.Port <= 0 || config.Port > 65535 {
			v.errors = append(v.errors, "email port must be between 1 and 65535")
		}

		if config.Username == "" {
			v.errors = append(v.errors, "email username is required when email host is provided")
		}

		if config.From == "" {
			v.errors = append(v.errors, "email from address is required when email host is provided")
		}
	}
}

// validateApp validates application configuration
func (v *Validator) validateApp(config AppConfig) {
	if config.Name == "" {
		v.errors = append(v.errors, "application name is required")
	}

	validEnvironments := []string{"development", "staging", "production", "test"}
	valid := false
	for _, env := range validEnvironments {
		if strings.ToLower(config.Environment) == env {
			valid = true
			break
		}
	}
	if !valid {
		v.errors = append(v.errors, fmt.Sprintf("application environment must be one of: %s", strings.Join(validEnvironments, ", ")))
	}

	if config.Version == "" {
		v.errors = append(v.errors, "application version is required")
	}
}

// ValidateConnectionString validates if a connection string is reachable
func (v *Validator) ValidateConnectionString(host, port string) error {
	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return fmt.Errorf("cannot connect to %s: %w", address, err)
	}
	defer conn.Close()
	return nil
}

// ValidatePort validates if a port is available
func (v *Validator) ValidatePort(port string) error {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("invalid port number: %s", port)
	}

	if portNum < 1 || portNum > 65535 {
		return fmt.Errorf("port number must be between 1 and 65535")
	}

	return nil
}
