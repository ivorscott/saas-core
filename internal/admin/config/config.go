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
		Port            string        `conf:"default:4000"`
	}
	Cognito struct {
		UserPoolID             string `conf:"required"`
		UserPoolClientID       string `conf:"required"`
		SharedUserPoolID       string `conf:"required"`
		SharedUserPoolClientID string `conf:"required"`
		M2MClientKey           string `conf:"required"`
		M2MClientSecret        string `conf:"required"`
		Region                 string `conf:"required"`
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
	Billing struct {
		ServiceAddress string `conf:"default:localhost"`
		ServicePort    string `conf:"default:4006"`
	}
}

// NewConfig returns a new Config.
func NewConfig() (Config, error) {
	var cfg Config

	if err := conf.Parse(os.Args[1:], "ADMIN", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			var usage string
			usage, err = conf.Usage("ADMIN", &cfg)
			if err != nil {
				panic(fmt.Errorf("error generating config usage %s", err.Error()))
			}
			println(usage)
			return cfg, err
		}
		panic(fmt.Errorf("error parsing config %s", err.Error()))
		return cfg, err
	}
	return cfg, nil
}
