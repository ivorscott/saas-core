// Package config manages configuration values.
package config

import (
	"fmt"
	"github.com/ardanlabs/conf"
	"os"
	"time"
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
		SharedUserPoolID string `conf:"required"`
		Region           string `conf:"required"`
	}
	DB struct {
		User       string `conf:"default:user_a,noprint"`
		Password   string `conf:"default:postgres,noprint"`
		Host       string `conf:"default:localhost,noprint"`
		Port       int    `conf:"default:5432,noprint"`
		Name       string `conf:"default:user,noprint"`
		DisableTLS bool   `conf:"default:false"`
	}
	Dynamodb struct {
		ConnectionTable string `conf:"required"`
	}

	Nats struct {
		Address string `conf:"default:127.0.0.1"`
		Port    string `conf:"default:4222"`
	}
}

// NewConfig returns a new Config.
func NewConfig() (Config, error) {
	var cfg Config

	if err := conf.Parse(os.Args[1:], "USER", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			var usage string
			usage, err = conf.Usage("USER", &cfg)
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
