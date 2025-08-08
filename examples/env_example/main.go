package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sublimeai21/config"
)

func main() {
	// Example 1: Load configuration from environment variables
	fmt.Println("=== Loading from Environment Variables ===")

	// Set environment variables (in real applications, these would be set by your system)
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("SERVER_HOST", "127.0.0.1")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_NAME", "example")
	os.Setenv("JWT_SECRET", "your-super-secret-jwt-key-that-is-at-least-32-characters-long")
	os.Setenv("APP_NAME", "Environment Example")
	os.Setenv("APP_ENVIRONMENT", "development")
	os.Setenv("APP_VERSION", "1.0.0")

	manager := config.NewManager()
	err := manager.Load(config.EnvironmentStrategy)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Access configuration
	serverConfig := manager.GetServerConfig()
	fmt.Printf("Server will run on %s:%s\n", serverConfig.Host, serverConfig.Port)

	dbConfig := manager.GetDatabaseConfig()
	fmt.Printf("Database: %s@%s:%s/%s\n", dbConfig.User, dbConfig.Host, dbConfig.Port, dbConfig.DBName)

	appConfig := manager.GetAppConfig()
	fmt.Printf("Application: %s (v%s) in %s mode\n",
		appConfig.Name, appConfig.Version, appConfig.Environment)

	// Example 2: Using helper methods with environment variables
	fmt.Println("\n=== Using Helper Methods ===")
	fmt.Printf("Database DSN: %s\n", manager.GetDatabaseDSN())
	fmt.Printf("Write Database DSN: %s\n", manager.GetWriteDatabaseDSN())
	fmt.Printf("Read Database DSN: %s\n", manager.GetReadDatabaseDSN())
	fmt.Printf("Is Read/Write Database: %t\n", manager.IsReadWriteDatabase())
	fmt.Printf("Database Config Type: %s\n", manager.GetDatabaseConfigType())
	fmt.Printf("Redis Address: %s\n", manager.GetRedisAddr())
	fmt.Printf("Server Address: %s\n", manager.GetServerAddr())
	fmt.Printf("Is Development: %t\n", manager.IsDevelopment())
	fmt.Printf("Is Production: %t\n", manager.IsProduction())
	fmt.Printf("Is Debug: %t\n", manager.IsDebug())

	// Example 3: Validation with environment variables
	fmt.Println("\n=== Configuration Validation ===")
	if err := manager.ValidateCurrent(); err != nil {
		fmt.Printf("Configuration validation failed: %v\n", err)
	} else {
		fmt.Println("Configuration is valid!")
	}

	// Example 4: Reloading configuration
	fmt.Println("\n=== Reloading Configuration ===")

	// Change an environment variable
	os.Setenv("SERVER_PORT", "8081")

	err = manager.Reload()
	if err != nil {
		log.Printf("Failed to reload configuration: %v", err)
	} else {
		serverConfig = manager.GetServerConfig()
		fmt.Printf("Server will now run on %s:%s\n", serverConfig.Host, serverConfig.Port)
	}

	// Example 5: Accessing specific configurations
	fmt.Println("\n=== Accessing Specific Configurations ===")

	jwtConfig := manager.GetJWTConfig()
	fmt.Printf("JWT Secret length: %d characters\n", len(jwtConfig.Secret))
	fmt.Printf("JWT Expiration: %v\n", jwtConfig.Expiration)

	logConfig := manager.GetLogConfig()
	fmt.Printf("Log Level: %s\n", logConfig.Level)
	fmt.Printf("Log Format: %s\n", logConfig.Format)

	emailConfig := manager.GetEmailConfig()
	if emailConfig.Host != "" {
		fmt.Printf("Email Host: %s:%d\n", emailConfig.Host, emailConfig.Port)
	} else {
		fmt.Println("Email configuration not provided")
	}
}

/*
To use this example with a .env file:

1. Create a .env file in your project root:
   ```
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

2. Load the .env file before running your application:
   ```bash
   # Using dotenv (if you have it installed)
   dotenv -f .env go run main.go

   # Or export the variables manually
   export $(cat .env | xargs)
   go run main.go
   ```

3. Your application will automatically pick up the environment variables.
*/
