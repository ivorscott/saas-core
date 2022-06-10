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
		UserPoolClientID string `conf:"required"`
		Region           string `conf:"default:eu-central-1"`
	}
	Nats struct {
		QueueGroup string `conf:"default:user"`
		Address    string `conf:"default:127.0.0.1"`
		Port       string `conf:"default:4222"`
	}
}
