# Go Configuration Shared Package

A comprehensive, type-safe configuration management package for Go applications that supports multiple configuration sources, validation, and runtime configuration changes.

## Features

- **Multiple Configuration Sources**: Environment variables, YAML files, and hybrid approaches
- **Type-Safe Configuration**: Strongly typed configuration structures
- **Validation**: Built-in configuration validation with custom rules
- **Thread-Safe**: Concurrent access to configuration with read-write locks
- **Configuration Watchers**: Subscribe to configuration changes
- **Helper Methods**: Convenient methods for common configuration access patterns
- **Environment-Aware**: Different strategies for different environments

## Installation

```bash
go get github.com/shared_pkg/config
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/shared_pkg/config"
)

func main() {
    // Create a configuration manager
    manager := config.NewManager()
    
    // Load configuration from environment variables
    err := manager.Load(config.EnvironmentStrategy)
    if err != nil {
        log.Fatal(err)
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

### Logging Configuration
```go
type LogConfig struct {
    Level      string `mapstructure:"level"`
    Format     string `mapstructure:"format"`
    OutputPath string `mapstructure:"output_path"`
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

### Application Configuration
```go
type AppConfig struct {
    Name        string `mapstructure:"name"`
    Environment string `mapstructure:"environment"`
    Version     string `mapstructure:"version"`
    Debug       bool   `mapstructure:"debug"`
}
```

## Configuration Sources

### Environment Variables

Set environment variables with the following naming convention:
- `SERVER_PORT` → Server.Port
- `DB_HOST` → Database.Host
- `JWT_SECRET` → JWT.Secret

Example:
```bash
export SERVER_PORT=8080
export DB_HOST=localhost
export DB_USER=postgres
export JWT_SECRET=your-super-secret-jwt-key-that-is-at-least-32-characters-long
```

### Configuration Files

Create a YAML configuration file:

```yaml
server:
  port: "8080"
  host: "0.0.0.0"
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "60s"

database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "password"
  dbname: "myapp"
  sslmode: "disable"
  max_conns: 10

redis:
  host: "localhost"
  port: "6379"
  password: ""
  db: 0

log:
  level: "info"
  format: "json"
  output_path: ""

jwt:
  secret: "your-super-secret-jwt-key-that-is-at-least-32-characters-long"
  expiration: "24h"
  issuer: "myapp"

email:
  host: "smtp.gmail.com"
  port: 587
  username: "your-email@gmail.com"
  password: "your-app-password"
  from: "noreply@myapp.com"

app:
  name: "My Application"
  environment: "development"
  version: "1.0.0"
  debug: true
```

## Usage Examples

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    "github.com/shared_pkg/config"
)

func main() {
    manager := config.NewManager()
    
    // Load from environment variables
    err := manager.Load(config.EnvironmentStrategy)
    if err != nil {
        log.Fatal(err)
    }
    
    // Access configuration
    serverConfig := manager.GetServerConfig()
    fmt.Printf("Server: %s:%s\n", serverConfig.Host, serverConfig.Port)
}
```

### Loading from File

```go
// Set the config file path
os.Setenv("CONFIG_PATH", "config.yaml")

manager := config.NewManager()
err := manager.Load(config.FileStrategy)
if err != nil {
    log.Fatal(err)
}
```

### Hybrid Loading (File with Environment Fallback)

```go
manager := config.NewManager()
err := manager.Load(config.HybridStrategy)
if err != nil {
    log.Fatal(err)
}
```

### Configuration Validation

```go
manager := config.NewManager()
err := manager.Load(config.EnvironmentStrategy)
if err != nil {
    log.Fatal(err)
}

// Validate the configuration
if err := manager.ValidateCurrent(); err != nil {
    log.Printf("Configuration validation failed: %v", err)
}
```

### Configuration Watchers

```go
type MyConfigWatcher struct{}

func (w *MyConfigWatcher) OnConfigChanged(oldConfig, newConfig *config.Config) {
    fmt.Println("Configuration changed!")
    if oldConfig != nil && newConfig != nil {
        fmt.Printf("Server port changed from %s to %s\n", 
            oldConfig.Server.Port, newConfig.Server.Port)
    }
}

// Add watcher
watcher := &MyConfigWatcher{}
manager.AddWatcher(watcher)
```

### Helper Methods

```go
// Get connection strings
dbDSN := manager.GetDatabaseDSN()
redisAddr := manager.GetRedisAddr()
serverAddr := manager.GetServerAddr()

// Environment checks
isDev := manager.IsDevelopment()
isProd := manager.IsProduction()
isDebug := manager.IsDebug()
```

## Configuration Strategies

### EnvironmentStrategy
Loads configuration from environment variables only.

### FileStrategy
Loads configuration from a configuration file (YAML, JSON, etc.).

### HybridStrategy
Attempts to load from file first, falls back to environment variables.

## Validation Rules

The package includes built-in validation for:

- **Server Configuration**: Port numbers, timeouts, host addresses
- **Database Configuration**: Required fields, SSL modes, connection limits
- **Redis Configuration**: Host, port, database numbers
- **Logging Configuration**: Valid log levels and formats
- **JWT Configuration**: Secret length, expiration times
- **Email Configuration**: Port ranges, required fields when host is provided
- **Application Configuration**: Environment values, version format

## Environment Variables Reference

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Server port | `8080` |
| `SERVER_HOST` | Server host | `0.0.0.0` |
| `SERVER_READ_TIMEOUT` | Read timeout | `30s` |
| `SERVER_WRITE_TIMEOUT` | Write timeout | `30s` |
| `SERVER_IDLE_TIMEOUT` | Idle timeout | `60s` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `` |
| `DB_NAME` | Database name | `app` |
| `DB_SSL_MODE` | SSL mode | `disable` |
| `DB_MAX_CONNS` | Max connections | `10` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `REDIS_PASSWORD` | Redis password | `` |
| `REDIS_DB` | Redis database | `0` |
| `LOG_LEVEL` | Log level | `info` |
| `LOG_FORMAT` | Log format | `json` |
| `LOG_OUTPUT_PATH` | Log output path | `` |
| `JWT_SECRET` | JWT secret | `your-secret-key` |
| `JWT_EXPIRATION` | JWT expiration | `24h` |
| `JWT_ISSUER` | JWT issuer | `app` |
| `EMAIL_HOST` | Email host | `` |
| `EMAIL_PORT` | Email port | `587` |
| `EMAIL_USERNAME` | Email username | `` |
| `EMAIL_PASSWORD` | Email password | `` |
| `EMAIL_FROM` | Email from address | `` |
| `APP_NAME` | Application name | `app` |
| `APP_ENVIRONMENT` | Environment | `development` |
| `APP_VERSION` | Application version | `1.0.0` |
| `APP_DEBUG` | Debug mode | `false` |
| `CONFIG_PATH` | Configuration file path | `config.yaml` |

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License.
