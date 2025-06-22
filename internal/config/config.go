package config

import (
	"net/url"
	"time"
)

// Load populates the application configuration from environment variables. If variables are not set, default
// values are applied instead. Load will panic if values from environment variables can't be converted to the
// desired data type.
func Load() Application {
	app := Application{
		ID:           envString("APPLICATION_ID", "Golang API Service"),
		BaseURL:      envURL("BASE_URL", "http://localhost:8080"),
		Address:      envString("ADDRESS", ":8080"),
		ReadTimeout:  envDuration("READ_TIMEOUT", 10*time.Second),
		WriteTimeout: envDuration("WRITE_TIMEOUT", 20*time.Second),
		IdleTimeout:  envDuration("IDLE_TIMEOUT", 30*time.Second),
		OAuth:        OAuth{},
		Database: Database{
			User:         envString("DB_USER", "postgres"),
			Password:     envString("DB_PASSWORD", "postgres"),
			Host:         envString("DB_HOST", "localhost"),
			Port:         envInt("DB_PORT", 5432),
			Name:         envString("DB_NAME", "golang_api_service"),
			SSL:          envBool("DB_SSLMODE", false),
			LogQueries:   envBool("DB_LOG_QUERIES", false),
			MigrateUp:    envBool("DB_MIGRATE_UP", true),
			LoadFixtures: envBool("DB_LOAD_FIXTURES", true),
		},
		HTTPClient: HTTPClient{},
	}
	return app
}

// Application represents the service main configuration.
type Application struct {
	ID           string
	BaseURL      *url.URL
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	OAuth        OAuth
	Database     Database
	HTTPClient   HTTPClient
}

// OAuth holds all configuration settings related to authorization.
type OAuth struct {
	InfoURL               *url.URL
	CredentialsDir        string
	ClientID              string
	EventBusReverserScope string
}

// Database holds values required to set up a db connection.
type Database struct {
	User         string
	Password     string
	Host         string
	Port         int
	Name         string
	SSL          bool
	LogQueries   bool
	MigrateUp    bool
	LoadFixtures bool
}

// HTTPClient holds values required to set up a API connection.
type HTTPClient struct {
	MinRetryInterval time.Duration
	MaxRetryInterval time.Duration
	RetryTimeout     time.Duration
	RequestTimeout   time.Duration
}

// TestCfg returns the configuration for testing purpose.
func (d Database) TestCfg() Database {
	cfg := d
	cfg.Port = envInt("DB_PORT", 5432)
	return cfg
}
