package config

import (
	"os"
	"testing"
	"time"

	"github.com/sublimeai21/config"
)

func TestNewManager(t *testing.T) {
	manager := config.NewManager()
	if manager == nil {
		t.Fatal("NewManager() returned nil")
	}

	if !manager.IsLoaded() {
		t.Log("Manager should not be loaded initially")
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	// Set up test environment variables
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("SERVER_HOST", "127.0.0.1")
	os.Setenv("DB_HOST", "test-db")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("JWT_SECRET", "test-secret-that-is-long-enough-for-validation")
	os.Setenv("APP_NAME", "Test App")
	os.Setenv("APP_ENVIRONMENT", "test")
	os.Setenv("APP_VERSION", "1.0.0")

	manager := config.NewManager()
	err := manager.Load(config.EnvironmentStrategy)
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	if !manager.IsLoaded() {
		t.Fatal("Manager should be loaded after successful load")
	}

	config := manager.GetConfig()
	if config == nil {
		t.Fatal("GetConfig() returned nil")
	}

	// Test server configuration
	serverConfig := manager.GetServerConfig()
	if serverConfig.Port != "9090" {
		t.Errorf("Expected server port 9090, got %s", serverConfig.Port)
	}
	if serverConfig.Host != "127.0.0.1" {
		t.Errorf("Expected server host 127.0.0.1, got %s", serverConfig.Host)
	}

	// Test database configuration
	dbConfig := manager.GetDatabaseConfig()
	if dbConfig.Host != "test-db" {
		t.Errorf("Expected database host test-db, got %s", dbConfig.Host)
	}
	if dbConfig.User != "testuser" {
		t.Errorf("Expected database user testuser, got %s", dbConfig.User)
	}
	if dbConfig.DBName != "testdb" {
		t.Errorf("Expected database name testdb, got %s", dbConfig.DBName)
	}

	// Test app configuration
	appConfig := manager.GetAppConfig()
	if appConfig.Name != "Test App" {
		t.Errorf("Expected app name 'Test App', got %s", appConfig.Name)
	}
	if appConfig.Environment != "test" {
		t.Errorf("Expected environment 'test', got %s", appConfig.Environment)
	}
}

func TestLoadFromFile(t *testing.T) {
	// Clear any existing environment variables that might interfere
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("SERVER_HOST")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("APP_NAME")
	os.Unsetenv("APP_ENVIRONMENT")
	os.Unsetenv("APP_VERSION")

	// Create a temporary config file
	configContent := `
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
  dbname: "testdb"
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
  secret: "test-secret-that-is-long-enough-for-validation"
  expiration: "24h"
  issuer: "testapp"

email:
  host: "smtp.test.com"
  port: 587
  username: "test@test.com"
  password: "password"
  from: "noreply@test.com"

app:
  name: "Test Application"
  environment: "test"
  version: "1.0.0"
  debug: true
`

	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(configContent)
	if err != nil {
		t.Fatalf("Failed to write config content: %v", err)
	}
	tmpFile.Close()

	// Set the config path
	os.Setenv("CONFIG_PATH", tmpFile.Name())

	manager := config.NewManager()
	err = manager.Load(config.FileStrategy)
	if err != nil {
		t.Fatalf("Failed to load configuration from file: %v", err)
	}

	config := manager.GetConfig()
	if config == nil {
		t.Fatal("GetConfig() returned nil")
	}

	// Test server configuration
	serverConfig := manager.GetServerConfig()
	if serverConfig.Port != "8080" {
		t.Errorf("Expected server port 8080, got %s", serverConfig.Port)
	}

	// Test database configuration
	dbConfig := manager.GetDatabaseConfig()
	if dbConfig.Host != "localhost" {
		t.Errorf("Expected database host localhost, got %s", dbConfig.Host)
	}
	if dbConfig.User != "postgres" {
		t.Errorf("Expected database user postgres, got %s", dbConfig.User)
	}

	// Test app configuration
	appConfig := manager.GetAppConfig()
	if appConfig.Name != "Test Application" {
		t.Errorf("Expected app name 'Test Application', got %s", appConfig.Name)
	}
	if !appConfig.Debug {
		t.Error("Expected debug to be true")
	}
}

func TestValidation(t *testing.T) {
	manager := config.NewManager()

	// Test with valid configuration
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("SERVER_HOST", "0.0.0.0")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("JWT_SECRET", "test-secret-that-is-long-enough-for-validation")
	os.Setenv("APP_NAME", "Test App")
	os.Setenv("APP_ENVIRONMENT", "development")
	os.Setenv("APP_VERSION", "1.0.0")

	err := manager.Load(config.EnvironmentStrategy)
	if err != nil {
		t.Fatalf("Failed to load valid configuration: %v", err)
	}

	// Validate should pass
	err = manager.ValidateCurrent()
	if err != nil {
		t.Errorf("Validation should pass for valid config: %v", err)
	}

	// Test with invalid configuration (missing required fields)
	manager2 := config.NewManager()
	os.Setenv("SERVER_PORT", "")     // Invalid: empty port
	os.Setenv("DB_HOST", "")         // Invalid: empty host
	os.Setenv("JWT_SECRET", "short") // Invalid: too short
	os.Setenv("APP_NAME", "Test App")
	os.Setenv("APP_ENVIRONMENT", "development")
	os.Setenv("APP_VERSION", "1.0.0")

	err = manager2.Load(config.EnvironmentStrategy)
	// The load should fail due to validation
	if err == nil {
		t.Error("Load should fail for invalid configuration")
	}

	// Check that the error contains validation information
	if err != nil {
		t.Logf("Expected validation error: %v", err)
	}
}

func TestHelperMethods(t *testing.T) {
	manager := config.NewManager()

	// Set up test configuration
	os.Setenv("SERVER_HOST", "localhost")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "pass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_SSL_MODE", "require")
	os.Setenv("REDIS_HOST", "redis.example.com")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("APP_ENVIRONMENT", "development")
	os.Setenv("APP_DEBUG", "true")
	os.Setenv("JWT_SECRET", "test-secret-that-is-long-enough-for-validation")
	os.Setenv("APP_NAME", "Test App")
	os.Setenv("APP_VERSION", "1.0.0")

	err := manager.Load(config.EnvironmentStrategy)
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Test helper methods
	expectedDSN := "host=db.example.com port=5432 user=user password=pass dbname=testdb sslmode=require"
	dsn := manager.GetDatabaseDSN()
	if dsn != expectedDSN {
		t.Errorf("Expected DSN '%s', got '%s'", expectedDSN, dsn)
	}

	expectedRedisAddr := "redis.example.com:6379"
	redisAddr := manager.GetRedisAddr()
	if redisAddr != expectedRedisAddr {
		t.Errorf("Expected Redis address '%s', got '%s'", expectedRedisAddr, redisAddr)
	}

	expectedServerAddr := "localhost:8080"
	serverAddr := manager.GetServerAddr()
	if serverAddr != expectedServerAddr {
		t.Errorf("Expected server address '%s', got '%s'", expectedServerAddr, serverAddr)
	}

	// Test environment checks
	if !manager.IsDevelopment() {
		t.Error("Expected IsDevelopment() to return true")
	}

	if manager.IsProduction() {
		t.Error("Expected IsProduction() to return false")
	}

	if !manager.IsDebug() {
		t.Error("Expected IsDebug() to return true")
	}
}

func TestConfigWatcher(t *testing.T) {
	manager := config.NewManager()

	// Create a test watcher
	watcherCalled := false
	var oldPort, newPort string

	watcher := &testConfigWatcher{
		onChanged: func(oldConfig, newConfig *config.Config) {
			watcherCalled = true
			if oldConfig != nil {
				oldPort = oldConfig.Server.Port
			}
			if newConfig != nil {
				newPort = newConfig.Server.Port
			}
		},
	}

	manager.AddWatcher(watcher)

	// Load initial configuration
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("APP_NAME", "Test App")
	os.Setenv("APP_ENVIRONMENT", "test")
	os.Setenv("APP_VERSION", "1.0.0")
	os.Setenv("JWT_SECRET", "test-secret-that-is-long-enough-for-validation")

	err := manager.Load(config.EnvironmentStrategy)
	if err != nil {
		t.Fatalf("Failed to load initial configuration: %v", err)
	}

	// Change configuration
	os.Setenv("SERVER_PORT", "9090")

	err = manager.Reload()
	if err != nil {
		t.Fatalf("Failed to reload configuration: %v", err)
	}

	// Wait a bit for the watcher to be called
	time.Sleep(100 * time.Millisecond)

	if !watcherCalled {
		t.Error("Config watcher was not called")
	}

	if oldPort != "8080" {
		t.Errorf("Expected old port 8080, got %s", oldPort)
	}

	if newPort != "9090" {
		t.Errorf("Expected new port 9090, got %s", newPort)
	}
}

// testConfigWatcher is a test implementation of ConfigWatcher
type testConfigWatcher struct {
	onChanged func(oldConfig, newConfig *config.Config)
}

func (w *testConfigWatcher) OnConfigChanged(oldConfig, newConfig *config.Config) {
	if w.onChanged != nil {
		w.onChanged(oldConfig, newConfig)
	}
}

// TestParseFunctions removed - parseInt and parseBool are private functions

func TestValidator(t *testing.T) {
	validator := config.NewValidator()

	// Test valid configuration
	validConfig := &config.Config{
		Server: config.ServerConfig{
			Port:         "8080",
			Host:         "0.0.0.0",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "password",
			DBName:   "testdb",
			SSLMode:  "disable",
			MaxConns: 10,
		},
		Redis: config.RedisConfig{
			Host:     "localhost",
			Port:     "6379",
			Password: "",
			DB:       0,
		},
		Log: config.LogConfig{
			Level:      "info",
			Format:     "json",
			OutputPath: "",
		},
		JWT: config.JWTConfig{
			Secret:     "test-secret-that-is-long-enough-for-validation",
			Expiration: 24 * time.Hour,
			Issuer:     "testapp",
		},
		App: config.AppConfig{
			Name:        "Test App",
			Environment: "development",
			Version:     "1.0.0",
			Debug:       false,
		},
	}

	err := validator.Validate(validConfig)
	if err != nil {
		t.Errorf("Validation should pass for valid config: %v", err)
	}

	// Test invalid configuration
	invalidConfig := &config.Config{
		Server: config.ServerConfig{
			Port: "", // Invalid: empty port
		},
		Database: config.DatabaseConfig{
			Host: "", // Invalid: empty host
		},
		JWT: config.JWTConfig{
			Secret: "short", // Invalid: too short
		},
	}

	err = validator.Validate(invalidConfig)
	if err == nil {
		t.Error("Validation should fail for invalid config")
	}

	validationErr, ok := err.(*config.ValidationError)
	if !ok {
		t.Error("Expected ValidationError type")
	}

	if len(validationErr.Errors) == 0 {
		t.Error("Expected validation errors")
	}
}
