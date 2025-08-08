package config

import (
	"fmt"
	"sync"
)

// Manager provides a high-level interface for configuration management
type Manager struct {
	config    *Config
	loader    *Loader
	validator *Validator
	mutex     sync.RWMutex
	watchers  []ConfigWatcher
}

// ConfigWatcher defines an interface for configuration change watchers
type ConfigWatcher interface {
	OnConfigChanged(oldConfig, newConfig *Config)
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	return &Manager{
		loader:    NewLoader(),
		validator: NewValidator(),
		watchers:  make([]ConfigWatcher, 0),
	}
}

// Load loads and validates configuration using the specified strategy
func (m *Manager) Load(strategy LoadStrategy) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	config, err := m.loader.Load(strategy)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate the configuration
	if err := m.validator.Validate(config); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Store the old config for watchers
	oldConfig := m.config
	m.config = config

	// Notify watchers if this is not the initial load
	if oldConfig != nil {
		m.notifyWatchers(oldConfig, config)
	}

	return nil
}

// GetConfig returns the current configuration (thread-safe)
func (m *Manager) GetConfig() *Config {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.config
}

// GetServerConfig returns the server configuration
func (m *Manager) GetServerConfig() ServerConfig {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if m.config == nil {
		return ServerConfig{}
	}
	return m.config.Server
}

// GetDatabaseConfig returns the database configuration
func (m *Manager) GetDatabaseConfig() DatabaseConfig {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if m.config == nil {
		return DatabaseConfig{}
	}
	return m.config.Database
}

// GetRedisConfig returns the Redis configuration
func (m *Manager) GetRedisConfig() RedisConfig {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if m.config == nil {
		return RedisConfig{}
	}
	return m.config.Redis
}

// GetLogConfig returns the logging configuration
func (m *Manager) GetLogConfig() LogConfig {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if m.config == nil {
		return LogConfig{}
	}
	return m.config.Log
}

// GetJWTConfig returns the JWT configuration
func (m *Manager) GetJWTConfig() JWTConfig {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if m.config == nil {
		return JWTConfig{}
	}
	return m.config.JWT
}

// GetEmailConfig returns the email configuration
func (m *Manager) GetEmailConfig() EmailConfig {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if m.config == nil {
		return EmailConfig{}
	}
	return m.config.Email
}

// GetAppConfig returns the application configuration
func (m *Manager) GetAppConfig() AppConfig {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if m.config == nil {
		return AppConfig{}
	}
	return m.config.App
}

// AddWatcher adds a configuration change watcher
func (m *Manager) AddWatcher(watcher ConfigWatcher) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.watchers = append(m.watchers, watcher)
}

// RemoveWatcher removes a configuration change watcher
func (m *Manager) RemoveWatcher(watcher ConfigWatcher) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i, w := range m.watchers {
		if w == watcher {
			m.watchers = append(m.watchers[:i], m.watchers[i+1:]...)
			break
		}
	}
}

// notifyWatchers notifies all watchers of configuration changes
func (m *Manager) notifyWatchers(oldConfig, newConfig *Config) {
	for _, watcher := range m.watchers {
		go func(w ConfigWatcher) {
			w.OnConfigChanged(oldConfig, newConfig)
		}(watcher)
	}
}

// Reload reloads the configuration from the current source
func (m *Manager) Reload() error {
	// Determine the current strategy based on environment
	strategy := EnvironmentStrategy
	if m.config != nil && m.config.App.Environment == "production" {
		strategy = FileStrategy
	}

	return m.Load(strategy)
}

// IsLoaded returns true if configuration has been loaded
func (m *Manager) IsLoaded() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.config != nil
}

// ValidateCurrent validates the current configuration
func (m *Manager) ValidateCurrent() error {
	m.mutex.RLock()
	config := m.config
	m.mutex.RUnlock()

	if config == nil {
		return fmt.Errorf("no configuration loaded")
	}

	return m.validator.Validate(config)
}

// GetDatabaseDSN returns the database connection string
func (m *Manager) GetDatabaseDSN() string {
	config := m.GetDatabaseConfig()
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
}

// GetRedisAddr returns the Redis address
func (m *Manager) GetRedisAddr() string {
	config := m.GetRedisConfig()
	return fmt.Sprintf("%s:%s", config.Host, config.Port)
}

// GetServerAddr returns the server address
func (m *Manager) GetServerAddr() string {
	config := m.GetServerConfig()
	return fmt.Sprintf("%s:%s", config.Host, config.Port)
}

// IsDevelopment returns true if the application is in development mode
func (m *Manager) IsDevelopment() bool {
	config := m.GetAppConfig()
	return config.Environment == "development"
}

// IsProduction returns true if the application is in production mode
func (m *Manager) IsProduction() bool {
	config := m.GetAppConfig()
	return config.Environment == "production"
}

// IsDebug returns true if debug mode is enabled
func (m *Manager) IsDebug() bool {
	config := m.GetAppConfig()
	return config.Debug
}
