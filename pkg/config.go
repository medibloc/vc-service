package pkg

import "time"

type Config struct {
	Debug        bool          `envconfig:"DEBUG" default:"false"`
	Port         int           `envconfig:"PORT" required:"true"`
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"10s"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"10s"`
	IdleTimeout  time.Duration `envconfig:"IDLE_TIMEOUT" default:"60s"`
}
