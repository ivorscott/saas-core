package config

import "time"

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
}
