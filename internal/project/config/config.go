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
		Port            string        `conf:"default:4001"`
	}
	Cognito struct {
		UserPoolID string `conf:"required"`
		Region     string `conf:"required"`
	}
	DB struct {
		User       string `conf:"default:user_a,noprint"`
		Password   string `conf:"default:postgres,noprint"`
		Host       string `conf:"default:localhost,noprint"`
		Port       int    `conf:"default:5432,noprint"`
		Name       string `conf:"default:project,noprint"`
		DisableTLS bool   `conf:"default:false"`
	}
	Nats struct {
		Address string `conf:"default:127.0.0.1"`
		Port    string `conf:"default:4222"`
	}
}
