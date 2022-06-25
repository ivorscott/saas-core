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
		UserPoolID string `conf:"required"`
		Region     string `conf:"required"`
	}
	Dynamodb struct {
		TenantTable string `conf:"required"`
		AuthTable   string `conf:"required"`
		ConfigTable string `conf:"required"`
		Host        string `conf:"default:localhost"`
		Port        string `conf:"default:8000"`
	}
	Nats struct {
		Address string `conf:"default:127.0.0.1"`
		Port    string `conf:"default:4222"`
	}
}
