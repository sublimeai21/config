package main

import (
	"fmt"
	"log"
	"os"

	core "github.com/shared_pkg/config/core"
)

// ExampleConfigWatcher demonstrates how to implement a config watcher
type ExampleConfigWatcher struct{}

func (w *ExampleConfigWatcher) OnConfigChanged(oldConfig, newConfig *core.Config) {
	fmt.Println("Configuration changed!")
	if oldConfig != nil && newConfig != nil {
		fmt.Printf("Server port changed from %s to %s\n",
			oldConfig.Server.Port, newConfig.Server.Port)
	}
}

func main() {
	// Create a new configuration manager
	manager := core.NewManager()

	// Add a configuration watcher
	watcher := &ExampleConfigWatcher{}
	manager.AddWatcher(watcher)

	// Example 1: Load configuration from environment variables
	fmt.Println("=== Loading from Environment Variables ===")

	// Set some environment variables for demonstration
	os.Setenv("APP_NAME", "Example App")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_NAME", "example")
	os.Setenv("JWT_SECRET", "your-super-secret-jwt-key-that-is-at-least-32-characters-long")

	err := manager.Load(core.EnvironmentStrategy)
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

	// Example 2: Load configuration from file
	fmt.Println("\n=== Loading from Configuration File ===")

	// Set the config path environment variable
	os.Setenv("CONFIG_PATH", "examples/config.yaml")

	err = manager.Load(core.FileStrategy)
	if err != nil {
		log.Printf("Failed to load from file (this is expected if config.yaml doesn't exist): %v", err)
	} else {
		serverConfig = manager.GetServerConfig()
		fmt.Printf("Server will run on %s:%s\n", serverConfig.Host, serverConfig.Port)
	}

	// Example 3: Using helper methods
	fmt.Println("\n=== Using Helper Methods ===")

	fmt.Printf("Database DSN: %s\n", manager.GetDatabaseDSN())
	fmt.Printf("Redis Address: %s\n", manager.GetRedisAddr())
	fmt.Printf("Server Address: %s\n", manager.GetServerAddr())
	fmt.Printf("Is Development: %t\n", manager.IsDevelopment())
	fmt.Printf("Is Production: %t\n", manager.IsProduction())
	fmt.Printf("Is Debug: %t\n", manager.IsDebug())

	// Example 4: Validation
	fmt.Println("\n=== Configuration Validation ===")

	if err := manager.ValidateCurrent(); err != nil {
		fmt.Printf("Configuration validation failed: %v\n", err)
	} else {
		fmt.Println("Configuration is valid!")
	}

	// Example 5: Reloading configuration
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

	// Example 6: Accessing specific configuration sections
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
