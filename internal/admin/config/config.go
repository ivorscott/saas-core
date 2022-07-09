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
		UserPoolID       string `conf:"required"`
		UserPoolClientID string `conf:"required"`
		SharedUserPoolID string `conf:"required"`
		Region           string `conf:"required"`
	}
	DB struct {
		User       string `conf:"default:postgres,noprint"`
		Password   string `conf:"default:postgres,noprint"`
		Host       string `conf:"default:localhost,noprint"`
		Port       int    `conf:"default:5432,noprint"`
		Name       string `conf:"default:admin,noprint"`
		DisableTLS bool   `conf:"default:false"`
	}
	Registration struct {
		ServiceAddress string `conf:"default:localhost"`
		ServicePort    string `conf:"default:4001"`
	}
	Tenant struct {
		ServiceAddress string `conf:"default:localhost"`
		ServicePort    string `conf:"default:4002"`
	}
}
