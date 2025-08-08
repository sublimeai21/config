package config

import "time"

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Email    EmailConfig    `mapstructure:"email"`
	App      AppConfig      `mapstructure:"app"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	Host         string        `mapstructure:"host"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	// --- Read/Write Database Configuration (Recommended) ---
	// These fields are used when DATABASE_CONFIG_TYPE=read_write
	DBWriteHost     string `mapstructure:"write_host"`
	DBWritePort     string `mapstructure:"write_port"`
	DBWriteUser     string `mapstructure:"write_user"`
	DBWritePassword string `mapstructure:"write_password"`
	DBWriteName     string `mapstructure:"write_dbname"`

	DBReadHost     string `mapstructure:"read_host"`
	DBReadPort     string `mapstructure:"read_port"`
	DBReadUser     string `mapstructure:"read_user"`
	DBReadPassword string `mapstructure:"read_password"`
	DBReadName     string `mapstructure:"read_dbname"`

	// --- Legacy Database Configuration (Backward Compatibility) ---
	// These fields are used when DATABASE_CONFIG_TYPE=legacy
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`

	// --- Database Type and Environment ---
	SSLMode            string `mapstructure:"sslmode"`
	MaxConns           int    `mapstructure:"max_conns"`
	DBType             string `mapstructure:"type"`
	Environment        string `mapstructure:"environment"`
	DatabaseConfigType string `mapstructure:"config_type"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	OutputPath string `mapstructure:"output_path"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
	Issuer     string        `mapstructure:"issuer"`
}

// EmailConfig holds email configuration
type EmailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
	Version     string `mapstructure:"version"`
	Debug       bool   `mapstructure:"debug"`
}
