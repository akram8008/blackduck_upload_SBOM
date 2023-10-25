package config

import (
	"log"
	"os"
)

type Config struct {
	BlackduckUrl   string
	BlackduckToken string
	ProjectName    string
	ProjectVersion string
}

func (c *Config) New() {
	c.BlackduckUrl = os.Getenv("BLACKDUCK_URL")
	c.BlackduckToken = os.Getenv("BLACKDUCK_TOKEN")
	if len(c.BlackduckToken) == 0 || len(c.BlackduckToken) == 0 {
		log.Fatalln("BLACKDUCK_URL and BLACKDUCK_TOKEN are required")
	}

	c.ProjectName = os.Getenv("PROJECT_NAME")
	c.ProjectVersion = os.Getenv("PROJECT_VERSION")

}
