package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config represents application configuration
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Arxiv     ArxivConfig     `mapstructure:"arxiv"`
	Cache     CacheConfig     `mapstructure:"cache"`
	CORS      CORSConfig      `mapstructure:"cors"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Driver       string        `mapstructure:"driver"`
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Username     string        `mapstructure:"username"`
	Password     string        `mapstructure:"password"`
	Database     string        `mapstructure:"database"`
	MaxOpenConns int           `mapstructure:"max_open_conns"`
	MaxIdleConns int           `mapstructure:"max_idle_conns"`
	MaxLifetime  time.Duration `mapstructure:"max_lifetime"`
}

// JWTConfig represents JWT configuration
type JWTConfig struct {
	Secret    string        `mapstructure:"secret"`
	ExpiresIn time.Duration `mapstructure:"expires_in"`
}

// ArxivConfig represents arXiv API configuration
type ArxivConfig struct {
	BaseURL    string        `mapstructure:"base_url"`
	Timeout    time.Duration `mapstructure:"timeout"`
	MaxRetries int           `mapstructure:"max_retries"`
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

// RateLimitConfig represents rate limiting configuration
type RateLimitConfig struct {
	Enabled  bool `mapstructure:"enabled"`
	Requests int  `mapstructure:"requests"` // requests per minute
	Burst    int  `mapstructure:"burst"`    // burst size
	PerIP    bool `mapstructure:"per_ip"`   // limit per IP or global
}

// Load loads configuration from file
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Set defaults (without JWT secret)
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Override with environment variables
	overrideWithEnvVars(&config)

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")

	// Database defaults
	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.database", "papertok")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.max_lifetime", "5m")

	// JWT defaults (only expires_in, no secret default)
	viper.SetDefault("jwt.expires_in", "168h") // 7 days

	viper.SetDefault("arxiv.base_url", "http://export.arxiv.org/api/query")
	viper.SetDefault("arxiv.timeout", "10s")
	viper.SetDefault("arxiv.max_retries", 3)

	viper.SetDefault("cache.enabled", true)
	viper.SetDefault("cache.ttl", "300s") // 5 minutes

	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:5173", "http://localhost:3000"})

	// Rate limiting defaults
	viper.SetDefault("rate_limit.enabled", true)
	viper.SetDefault("rate_limit.requests", 60) // 60 requests per minute
	viper.SetDefault("rate_limit.burst", 10)
	viper.SetDefault("rate_limit.per_ip", true)
}

// overrideWithEnvVars overrides configuration with environment variables
func overrideWithEnvVars(config *Config) {
	// JWT Secret (required)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic(errors.New("JWT_SECRET environment variable is required"))
	}
	// Reject common insecure default values
	if jwtSecret == "change-this-secret-in-production" ||
		jwtSecret == "your-secret-key-change-in-production" ||
		jwtSecret == "your-secret-key" ||
		jwtSecret == "secret" {
		panic(errors.New("JWT_SECRET must be set to a secure value, not a default placeholder"))
	}
	config.JWT.Secret = jwtSecret

	// Server Configuration
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Server.Port = p
		}
	}
	if mode := os.Getenv("SERVER_MODE"); mode != "" {
		config.Server.Mode = mode
	}

	// Database Configuration
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Database.Port = p
		}
	}
	if username := os.Getenv("DB_USERNAME"); username != "" {
		config.Database.Username = username
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if database := os.Getenv("DB_DATABASE"); database != "" {
		config.Database.Database = database
	}

	// ArXiv Configuration
	if baseURL := os.Getenv("ARXIV_BASE_URL"); baseURL != "" {
		config.Arxiv.BaseURL = baseURL
	}
	if timeout := os.Getenv("ARXIV_TIMEOUT"); timeout != "" {
		if t, err := time.ParseDuration(timeout); err == nil {
			config.Arxiv.Timeout = t
		}
	}
	if retries := os.Getenv("ARXIV_MAX_RETRIES"); retries != "" {
		if r, err := strconv.Atoi(retries); err == nil {
			config.Arxiv.MaxRetries = r
		}
	}

	// Cache Configuration
	if enabled := os.Getenv("CACHE_ENABLED"); enabled != "" {
		config.Cache.Enabled = strings.ToLower(enabled) == "true"
	}
	if ttl := os.Getenv("CACHE_TTL"); ttl != "" {
		if t, err := time.ParseDuration(ttl); err == nil {
			config.Cache.TTL = t
		}
	}

	// CORS Configuration
	if origins := os.Getenv("CORS_ALLOWED_ORIGINS"); origins != "" {
		config.CORS.AllowedOrigins = strings.Split(origins, ",")
	}
}
