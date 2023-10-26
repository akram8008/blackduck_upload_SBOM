package main

import (
	"github.com/akram8008/blackduck_upload_SBOM/api/docs"
	"github.com/akram8008/blackduck_upload_SBOM/internal/app"
	"github.com/akram8008/blackduck_upload_SBOM/internal/config"
	log "github.com/sirupsen/logrus"
)

// @title Blackduck Api request
// @contact.name API Support
// @contact.urlhttps://github.com/akram8008
// @contact.email akram8008@gmail.com

func main() {
	l := log.New()
	var configEnv config.Config
	if err := config.New(&configEnv); err != nil {
		l.Fatalln("config: ", err)
	}
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http"}

	server := app.New(&configEnv, l)
	server.Run()
}
