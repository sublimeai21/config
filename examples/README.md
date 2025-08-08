# Configuration Package Examples

This directory contains examples demonstrating how to use the Go Configuration Management Package.

## Files

- `main.go` - Complete example showing all features
- `env_example.go` - Environment variables example
- `config.yaml` - YAML configuration file example
- `.env` - Environment variables file example

## Running Examples

### Basic Example
```bash
go run main.go
```

### Environment Variables Example
```bash
go run env_example.go
```

## Using Environment Variables

### Method 1: Direct Environment Variables
Set environment variables in your shell:
```bash
export SERVER_PORT=8080
export DB_HOST=localhost
export JWT_SECRET=your-secret-key
go run main.go
```

### Method 2: Using .env File
1. Create a `.env` file in your project root:
```env
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=myapp
JWT_SECRET=your-super-secret-jwt-key-that-is-at-least-32-characters-long
APP_NAME=My Application
APP_ENVIRONMENT=development
APP_VERSION=1.0.0
```

2. Load the environment variables:
```bash
# Using dotenv (if installed)
dotenv -f .env go run main.go

# Or export manually
export $(cat .env | xargs)
go run main.go
```

### Method 3: Using YAML File
1. Create a `config.yaml` file:
```yaml
server:
  port: "8080"
  host: "0.0.0.0"

database:
  host: "localhost"
  user: "postgres"
  password: "password"
  dbname: "myapp"

jwt:
  secret: "your-super-secret-jwt-key-that-is-at-least-32-characters-long"
  expiration: "24h"
```

2. Set the config path and run:
```bash
export CONFIG_PATH=config.yaml
go run main.go
```

## Environment Variables Reference

### Server Configuration
- `SERVER_PORT` - Server port (default: "8080")
- `SERVER_HOST` - Server host (default: "0.0.0.0")
- `SERVER_READ_TIMEOUT` - Read timeout (default: "30s")
- `SERVER_WRITE_TIMEOUT` - Write timeout (default: "30s")
- `SERVER_IDLE_TIMEOUT` - Idle timeout (default: "60s")

### Database Configuration
- `DB_HOST` - Database host (default: "localhost")
- `DB_PORT` - Database port (default: "5432")
- `DB_USER` - Database user (default: "postgres")
- `DB_PASSWORD` - Database password (default: "")
- `DB_NAME` - Database name (default: "app")
- `DB_SSL_MODE` - SSL mode (default: "disable")
- `DB_MAX_CONNS` - Max connections (default: 10)

### Redis Configuration
- `REDIS_HOST` - Redis host (default: "localhost")
- `REDIS_PORT` - Redis port (default: "6379")
- `REDIS_PASSWORD` - Redis password (default: "")
- `REDIS_DB` - Redis database (default: 0)

### JWT Configuration
- `JWT_SECRET` - JWT secret (required, min 32 chars)
- `JWT_EXPIRATION` - JWT expiration (default: "24h")
- `JWT_ISSUER` - JWT issuer (default: "app")

### Logging Configuration
- `LOG_LEVEL` - Log level (default: "info")
- `LOG_FORMAT` - Log format (default: "json")
- `LOG_OUTPUT_PATH` - Log output path (default: "")

### Email Configuration
- `EMAIL_HOST` - Email host (default: "")
- `EMAIL_PORT` - Email port (default: 587)
- `EMAIL_USERNAME` - Email username (default: "")
- `EMAIL_PASSWORD` - Email password (default: "")
- `EMAIL_FROM` - Email from address (default: "")

### Application Configuration
- `APP_NAME` - Application name (default: "app")
- `APP_ENVIRONMENT` - Environment (default: "development")
- `APP_VERSION` - Version (default: "1.0.0")
- `APP_DEBUG` - Debug mode (default: false)

## Best Practices

1. **Never commit sensitive data** - Use environment variables for secrets
2. **Use .env for development** - Keep a `.env.example` file in your repo
3. **Validate configuration** - Always validate your config before using it
4. **Use appropriate defaults** - Set sensible defaults for all configuration
5. **Environment-specific configs** - Use different configs for dev/staging/prod

## Security Notes

- The `.env` file contains example values and should not be used in production
- Always use strong, unique secrets in production
- Consider using a secrets management service for production environments
- Never commit real passwords or API keys to version control
