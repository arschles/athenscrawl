package config

import "fmt"

type Config struct {
	Endpoint         string `envconfig:"GOPROXY" required:"true"`
	GHClientID       string `envconfig:"GITHUB_CLIENT_ID" required:"true"`
	GHClientSecret   string `envconfig:"GITHUB_CLIENT_SECRET"`
	GHTickDurSec     int    `envconfig:"GITHUB_API_TICK" default:"1"`
	AthensTickDurSec int    `envconfig:"ATHENS_TICK" default:"1"`
}

func (c *Config) String() string {
	return fmt.Sprintf(`GOPROXY: %s
GitHub Client ID: %s
GitHub Client Secret: <secure>
GitHub Tick Duration: %d sec
Athens Tick Duration: %d sec`,
		c.Endpoint,
		c.GHClientID,
		c.GHTickDurSec,
		c.AthensTickDurSec,
	)
}
