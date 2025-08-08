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
	Port         string        `mapstructure:"port"`          // e.g., "8080", "3000", "9090"
	Host         string        `mapstructure:"host"`          // e.g., "localhost", "0.0.0.0", "127.0.0.1"
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`  // e.g., "30s", "1m", "5m"
	WriteTimeout time.Duration `mapstructure:"write_timeout"` // e.g., "30s", "1m", "5m"
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`  // e.g., "60s", "2m", "10m"
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	// --- Read/Write Database Configuration (Recommended) ---
	// These fields are used when DATABASE_CONFIG_TYPE=read_write
	DBWriteHost     string `mapstructure:"write_host"`     // e.g., "write-db.example.com", "master-db.internal"
	DBWritePort     string `mapstructure:"write_port"`     // e.g., "5432", "3306", "1433"
	DBWriteUser     string `mapstructure:"write_user"`     // e.g., "write_user", "master_user"
	DBWritePassword string `mapstructure:"write_password"` // e.g., "write_password", "master_password"
	DBWriteName     string `mapstructure:"write_dbname"`   // e.g., "myapp_write", "master_db"

	DBReadHost     string `mapstructure:"read_host"`     // e.g., "read-db.example.com", "replica-db.internal"
	DBReadPort     string `mapstructure:"read_port"`     // e.g., "5432", "3306", "1433"
	DBReadUser     string `mapstructure:"read_user"`     // e.g., "read_user", "replica_user"
	DBReadPassword string `mapstructure:"read_password"` // e.g., "read_password", "replica_password"
	DBReadName     string `mapstructure:"read_dbname"`   // e.g., "myapp_read", "replica_db"

	// --- Legacy Database Configuration (Backward Compatibility) ---
	// These fields are used when DATABASE_CONFIG_TYPE=legacy
	Host     string `mapstructure:"host"`     // e.g., "localhost", "db.example.com", "127.0.0.1"
	Port     string `mapstructure:"port"`     // e.g., "5432", "3306", "1433"
	User     string `mapstructure:"user"`     // e.g., "postgres", "mysql_user", "sa"
	Password string `mapstructure:"password"` // e.g., "password", "secret", ""
	DBName   string `mapstructure:"dbname"`   // e.g., "myapp", "testdb", "production"

	// --- Database Type and Environment ---
	SSLMode            string `mapstructure:"sslmode"`     // e.g., "disable", "require", "verify-ca", "verify-full"
	MaxConns           int    `mapstructure:"max_conns"`   // e.g., 10, 50, 100
	DBType             string `mapstructure:"type"`        // e.g., "postgresql", "mysql", "sqlserver", "sqlite"
	Environment        string `mapstructure:"environment"` // e.g., "development", "staging", "production"
	DatabaseConfigType string `mapstructure:"config_type"` // e.g., "read_write", "legacy", "auto_detect"
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`     // e.g., "localhost", "redis.example.com", "127.0.0.1"
	Port     string `mapstructure:"port"`     // e.g., "6379", "6380", "26379"
	Password string `mapstructure:"password"` // e.g., "redis_password", "secret", ""
	DB       int    `mapstructure:"db"`       // e.g., 0, 1, 2, 15
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level      string `mapstructure:"level"`       // e.g., "debug", "info", "warn", "error", "fatal"
	Format     string `mapstructure:"format"`      // e.g., "json", "text", "logfmt"
	OutputPath string `mapstructure:"output_path"` // e.g., "/var/log/app.log", "stdout", "stderr"
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`     // e.g., "your-super-secret-jwt-key-here"
	Expiration time.Duration `mapstructure:"expiration"` // e.g., "24h", "7d", "30m"
	Issuer     string        `mapstructure:"issuer"`     // e.g., "myapp", "auth-service", "api-gateway"
}

// EmailConfig holds email configuration
type EmailConfig struct {
	Host     string `mapstructure:"host"`     // e.g., "smtp.gmail.com", "smtp.sendgrid.net", "mail.example.com"
	Port     int    `mapstructure:"port"`     // e.g., 587, 465, 25
	Username string `mapstructure:"username"` // e.g., "user@example.com", "noreply@myapp.com"
	Password string `mapstructure:"password"` // e.g., "email_password", "app_password"
	From     string `mapstructure:"from"`     // e.g., "noreply@myapp.com", "support@example.com"
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name        string `mapstructure:"name"`        // e.g., "My Application", "API Gateway", "User Service"
	Environment string `mapstructure:"environment"` // e.g., "development", "staging", "production", "test"
	Version     string `mapstructure:"version"`     // e.g., "1.0.0", "v2.1.3", "dev"
	Debug       bool   `mapstructure:"debug"`       // e.g., true, false
}
