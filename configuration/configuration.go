package configuration

import (
	"fmt"
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
	} `yaml:"session"`

	Database struct {
		URL      string `yaml:"url"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DbName   string `yaml:"dbName"`
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
	if c.Session.SessionKey == "" {
		return fmt.Errorf("incomplete session configuration: missing SessionKey")
	}
	if c.Session.SessionContextKey == "" {
		return fmt.Errorf("incomplete session configuration: missing SessionContextKey")
	}
	return nil
}
