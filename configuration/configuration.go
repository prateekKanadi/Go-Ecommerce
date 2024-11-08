package configuration

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// variable declaration
var (
	// global
	Conf *Config

	// local

)

type Config struct {
	Session struct {
		SessionKey        string `yaml:"session_key"`
		SessionContextKey string `yaml:"session_context_key"`
	} `yaml:"session"`

	Database struct {
		URL      string `yaml:"url"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DbName   string `yaml:"dbName"`
	} `yaml:"database"`
}

func Init(configPath string) *Config {
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	Conf = config
	return config
}

// helper functions
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
