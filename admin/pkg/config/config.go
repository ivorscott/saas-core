package config

import "time"

// Config contains application configuration with good defaults.
type Config struct {
	Web struct {
		Debug           string        `conf:"default:localhost:6060"`
		Production      bool          `conf:"default:false"`
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:5s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
		AppFrontend     string        `conf:"default:localhost:4000"`
		AppBackend      string        `conf:"default:http://localhost:4001"`
	}
}
