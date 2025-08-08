# Go Configuration Management Package

A comprehensive configuration management system for Go applications with support for multiple sources, validation, and thread-safe access.

## Features

- 🔧 **Multiple Configuration Sources**: Environment variables, YAML files, and hybrid loading
- 🛡️ **Validation**: Built-in configuration validation with customizable rules
- 🔄 **Hot Reloading**: Support for configuration changes at runtime
- 👀 **Watchers**: Event-driven configuration change notifications
- 🧵 **Thread-Safe**: Concurrent access with read-write locks
- 🎯 **Type-Safe**: Strongly typed configuration structures
- 🚀 **Easy to Use**: Simple API with helper methods

## Installation

```bash
go get github.com/sublimeai21/config/core
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    core "github.com/sublimeai21/config/core"
)

func main() {
    // Create a new configuration manager
    manager := core.NewManager()
    
    // Load configuration from environment variables
    err := manager.Load(core.EnvironmentStrategy)
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
    
    // Access configuration
    serverConfig := manager.GetServerConfig()
    fmt.Printf("Server: %s:%s\n", serverConfig.Host, serverConfig.Port)
    
    dbConfig := manager.GetDatabaseConfig()
    fmt.Printf("Database: %s@%s:%s/%s\n", 
        dbConfig.User, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
}
```

## Configuration Structure

The package supports the following configuration sections:

### Server Configuration
```go
type ServerConfig struct {
    Port         string        `mapstructure:"port"`
    Host         string        `mapstructure:"host"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
    IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}
```

### Database Configuration
```go
type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     string `mapstructure:"port"`
    User     string `mapstructure:"user"`
    Password string `mapstructure:"password"`
    DBName   string `mapstructure:"dbname"`
    SSLMode  string `mapstructure:"sslmode"`
    MaxConns int    `mapstructure:"max_conns"`
}
```

### Redis Configuration
```go
type RedisConfig struct {
    Host     string `mapstructure:"host"`
    Port     string `mapstructure:"port"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
}
```

### JWT Configuration
```go
type JWTConfig struct {
    Secret     string        `mapstructure:"secret"`
    Expiration time.Duration `mapstructure:"expiration"`
    Issuer     string        `mapstructure:"issuer"`
}
```

### Log Configuration
```go
type LogConfig struct {
    Level      string `mapstructure:"level"`
    Format     string `mapstructure:"format"`
    OutputPath string `mapstructure:"output_path"`
}
```

### Email Configuration
```go
type EmailConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    From     string `mapstructure:"from"`
}
```

### App Configuration
```go
type AppConfig struct {
    Name        string `mapstructure:"name"`
    Environment string `mapstructure:"environment"`
    Version     string `mapstructure:"version"`
    Debug       bool   `mapstructure:"debug"`
}
```

## Loading Strategies

### Environment Variables
```go
err := manager.Load(core.EnvironmentStrategy)
```

Environment variables are automatically mapped:
- `SERVER_PORT` → `Server.Port`
- `DB_HOST` → `Database.Host`
- `JWT_SECRET` → `JWT.Secret`
- etc.

### File-based Configuration
```go
os.Setenv("CONFIG_PATH", "config.yaml")
err := manager.Load(core.FileStrategy)
```

### Hybrid Strategy
```go
err := manager.Load(core.HybridStrategy)
```
Tries file first, falls back to environment variables.

## Configuration Watchers

Implement the `ConfigWatcher` interface to receive notifications when configuration changes:

```go
type MyWatcher struct{}

func (w *MyWatcher) OnConfigChanged(oldConfig, newConfig *core.Config) {
    fmt.Println("Configuration changed!")
    if oldConfig != nil && newConfig != nil {
        fmt.Printf("Server port changed from %s to %s\n",
            oldConfig.Server.Port, newConfig.Server.Port)
    }
}

// Add the watcher
watcher := &MyWatcher{}
manager.AddWatcher(watcher)
```

## Helper Methods

The manager provides convenient helper methods:

```go
// Connection strings
dsn := manager.GetDatabaseDSN()
redisAddr := manager.GetRedisAddr()
serverAddr := manager.GetServerAddr()

// Environment checks
isDev := manager.IsDevelopment()
isProd := manager.IsProduction()
isDebug := manager.IsDebug()

// Validation
err := manager.ValidateCurrent()

// Reloading
err := manager.Reload()
```

## Environment Variables

The package supports the following environment variables:

### Server
- `SERVER_PORT` (default: "8080")
- `SERVER_HOST` (default: "0.0.0.0")
- `SERVER_READ_TIMEOUT` (default: "30s")
- `SERVER_WRITE_TIMEOUT` (default: "30s")
- `SERVER_IDLE_TIMEOUT` (default: "60s")

### Database
- `DB_HOST` (default: "localhost")
- `DB_PORT` (default: "5432")
- `DB_USER` (default: "postgres")
- `DB_PASSWORD` (default: "")
- `DB_NAME` (default: "app")
- `DB_SSL_MODE` (default: "disable")
- `DB_MAX_CONNS` (default: 10)

### Redis
- `REDIS_HOST` (default: "localhost")
- `REDIS_PORT` (default: "6379")
- `REDIS_PASSWORD` (default: "")
- `REDIS_DB` (default: 0)

### JWT
- `JWT_SECRET` (default: "your-secret-key")
- `JWT_EXPIRATION` (default: "24h")
- `JWT_ISSUER` (default: "app")

### Logging
- `LOG_LEVEL` (default: "info")
- `LOG_FORMAT` (default: "json")
- `LOG_OUTPUT_PATH` (default: "")

### Email
- `EMAIL_HOST` (default: "")
- `EMAIL_PORT` (default: 587)
- `EMAIL_USERNAME` (default: "")
- `EMAIL_PASSWORD` (default: "")
- `EMAIL_FROM` (default: "")

### Application
- `APP_NAME` (default: "app")
- `APP_ENVIRONMENT` (default: "development")
- `APP_VERSION` (default: "1.0.0")
- `APP_DEBUG` (default: false)

## Validation

The package includes built-in validation:

- JWT secret must be at least 32 characters long
- JWT secret is required
- Custom validation rules can be added

## Examples

See the `examples/` directory for complete usage examples.

## Testing

Run the tests:

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
