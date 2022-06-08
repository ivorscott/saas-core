package config

import "time"

// Config represents the application configuration.
type Config struct {
	Web struct {
		Production      bool          `conf:"default:false"`
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:5s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
		Address         string        `conf:"default:localhost"`
		Port            string        `conf:"default:4000"`
	}
	Cognito struct {
		AppClientID      string `conf:"required"`
		UserPoolClientID string `conf:"required"`
		Region           string `conf:"default:eu-central-1"`
	}
	AWS struct {
		Region string `conf:"default:eu-central-1"`
	}
	Dynamodb struct {
		TenantTable string `conf:"required"`
		AuthTable   string `conf:"required"`
		ConfigTable string `conf:"required"`
		Port        string
	}
}
