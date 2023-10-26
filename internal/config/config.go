package config

import (
	"github.com/pkg/errors"
	"os"
)

type Config struct {
	BlackduckUrl   string
	BlackduckToken string
}

func New(c *Config) error {
	c.BlackduckUrl = os.Getenv("BLACKDUCK_URL")
	c.BlackduckToken = os.Getenv("BLACKDUCK_TOKEN")
	if c.BlackduckUrl == "" || c.BlackduckToken == "" {
		return errors.New("not all envs are provided")
	}
	return nil
}
