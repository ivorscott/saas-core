package config

import (
	"time"
)

// Config contains application configuration with good defaults.
type Config struct {
	Web struct {
		DebugPort       string        `conf:"default:6060"`
		Production      bool          `conf:"default:false"`
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:5s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
		Backend         string        `conf:"default:localhost:4001"`
		BackendPort     string        `conf:"default:4001"`
		FrontendPort    string        `conf:"default:4000"`
	}
	Cognito struct {
		AppClientID      string `conf:"required"`
		UserPoolClientID string `conf:"required"`
		Region           string `conf:"default:eu-central-1"`
	}
	Postgres struct {
		User       string `conf:"required"`
		Password   string `conf:"required"`
		Host       string `conf:"required"`
		Port       int    `conf:"required"`
		DB         string `conf:"required"`
		DisableTLS bool   `conf:"default:false"`
	}
}
