package config

import (
	"fmt"
	"time"
)

// Config is the configuration values from the environment
type Config struct {
	Endpoint        string `envconfig:"GOPROXY" required:"true"`
	GHClientID      string `envconfig:"GITHUB_CLIENT_ID" required:"true"`
	GHClientSecret  string `envconfig:"GITHUB_CLIENT_SECRET"`
	GHTickDurMS     int    `envconfig:"GITHUB_API_TICK" default:"1000"`
	AthensTickDurMS int    `envconfig:"ATHENS_TICK" default:"1000"`
	Debug           bool   `envconfig:"DEBUG" default:"false"`
}

func (c *Config) GHTickDur() time.Duration {
	return time.Duration(c.GHTickDurMS) * time.Millisecond
}

func (c *Config) AthensTickDur() time.Duration {
	return time.Duration(c.AthensTickDurMS) * time.Millisecond
}

func (c *Config) String() string {
	return fmt.Sprintf(`GOPROXY: %s
GitHub Client ID: %s
GitHub Client Secret: <secure>
GitHub Tick Duration: %d sec
Athens Tick Duration: %d sec`,
		c.Endpoint,
		c.GHClientID,
		c.GHTickDurMS,
		c.AthensTickDurMS,
	)
}
