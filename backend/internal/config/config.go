package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config represents application configuration
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Arxiv  ArxivConfig  `mapstructure:"arxiv"`
	Cache  CacheConfig  `mapstructure:"cache"`
	CORS   CORSConfig   `mapstructure:"cors"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release
}

// ArxivConfig represents arXiv API configuration
type ArxivConfig struct {
	BaseURL   string        `mapstructure:"base_url"`
	Timeout   time.Duration `mapstructure:"timeout"`
	MaxRetries int          `mapstructure:"max_retries"`
}

// CacheConfig represents cache configuration
type CacheConfig struct {
	Enabled bool          `mapstructure:"enabled"`
	TTL     time.Duration `mapstructure:"ttl"`
}

// CORSConfig represents CORS configuration
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

// Load loads configuration from file
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Set defaults
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("arxiv.base_url", "http://export.arxiv.org/api/query")
	viper.SetDefault("arxiv.timeout", "10s")
	viper.SetDefault("arxiv.max_retries", 3)
	viper.SetDefault("cache.enabled", true)
	viper.SetDefault("cache.ttl", "300s") // 5 minutes
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:5173", "http://localhost:3000"})
}
