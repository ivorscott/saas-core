// Package config manages configuration values.
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ardanlabs/conf"
)

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
	Dynamodb struct {
		AuthTable string `conf:"required"`
	}
	Nats struct {
		Address string `conf:"default:127.0.0.1"`
		Port    string `conf:"default:4222"`
	}
}

// NewConfig creates a new configuration struct for the service.
func NewConfig() (Config, error) {
	var cfg Config

	if err := conf.Parse(os.Args[1:], "REGISTRATION", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			var usage string
			usage, err = conf.Usage("REGISTRATION", &cfg)
			if err != nil {
				panic(fmt.Errorf("error generating config usage: %s", err.Error()))
			}
			println(usage)
			return cfg, err
		}
		panic(fmt.Errorf("error parsing config: %s", err.Error()))
		return cfg, err
	}
	return cfg, nil
}
