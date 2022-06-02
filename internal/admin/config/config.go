package config

import "time"

type Config struct {
	Web struct {
		DebugPort       string        `conf:"default:6060"`
		Production      bool          `conf:"default:false"`
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:5s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
		Address         string        `conf:"default:localhost"`
		Port            string        `conf:"default:4001"`
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
