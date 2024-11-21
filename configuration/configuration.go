package configuration

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	Session struct {
		SessionKey        string `yaml:"session_key"`
		SessionContextKey string `yaml:"session_context_key"`
		Domain            string `yaml:"domain"`
		Secure            bool   `yaml:"secure"`
		HttpOnly          bool   `yaml:"http_only"`
		Path              string `yaml:"path"`
		MaxAge            int    `yaml:"max_age"`
	} `yaml:"session"`

	Database struct {
		URL             string `yaml:"url"`
		User            string `yaml:"user"`
		Password        string `yaml:"password"`
		DbName          string `yaml:"dbName"`
		MaxOpenConns    int    `yaml:"max_open_conns"`
		MaxIdleConns    int    `yaml:"max_idle_conns"`
		ConnMaxLifetime int    `yaml:"conn_max_lifetime"` // In seconds
	} `yaml:"database"`
}

// Init loads and initializes the configuration from the specified file
func Init(configPath string) (*Config, error) {
	config, err := loadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return config, nil
}

// loadConfig loads configuration from a YAML file
func loadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Validate checks if the loaded configuration is complete and valid
func (c *Config) Validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("incomplete database configuration: missing URL")
	}
	if c.Database.User == "" {
		return fmt.Errorf("incomplete database configuration: missing User")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("incomplete database configuration: missing Password")
	}
	if c.Database.DbName == "" {
		return fmt.Errorf("incomplete database configuration: missing DbName")
	}
	if c.Database.MaxOpenConns <= 0 {
		return fmt.Errorf("database configuration error: MaxOpenConns must be greater than 0")
	}

	if c.Database.MaxIdleConns < 0 {
		return fmt.Errorf("database configuration error: MaxIdleConns cannot be negative")
	}

	if c.Database.ConnMaxLifetime < 0 {
		return fmt.Errorf("database configuration error: ConnMaxLifetime cannot be negative")
	}

	if c.Session.SessionKey == "" {
		return fmt.Errorf("incomplete session configuration: missing SessionKey")
	}
	if c.Session.SessionContextKey == "" {
		return fmt.Errorf("incomplete session configuration: missing SessionContextKey")
	}
	if c.Session.Domain == "" {
		return fmt.Errorf("incomplete session configuration: missing Domain")
	}
	// Validate Secure flag based on the Domain
	if c.Session.Domain == "localhost" && c.Session.Secure {
		log.Println("Warning: Secure is true for localhost. This may not work in development.")
	} else if c.Session.Domain != "localhost" && !c.Session.Secure {
		return errors.New("session configuration error: Secure must be true for production environments")
	}
	if c.Session.Path == "" {
		return errors.New("session configuration error: Path cannot be empty")
	}

	if c.Session.MaxAge < 0 {
		return errors.New("session configuration error: MaxAge cannot be negative")
	}

	return nil
}
