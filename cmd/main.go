package main

import (
	"context"
	"os"

	"github.com/akram8008/blackduck_upload_SBOM/pkg/blackduck"
	"github.com/akram8008/blackduck_upload_SBOM/pkg/config"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func main() {
	var configEnv config.Config
	configEnv.New()
	b := blackduck.NewBlackDuckComponents(configEnv.BlackduckUrl, configEnv.BlackduckToken)

	testIt(configEnv, b)

}

func testIt(configEnv config.Config, b *blackduck.BlackDuckComponents) {
	f, err := os.Open("./file/bom.xml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	req := blackduck.RequestUploadSBOM{
		File:           f,
		FileName:       "bom.xml",
		ProjectName:    configEnv.ProjectName,
		ProjectVersion: configEnv.ProjectVersion,
	}

	ctx := context.WithValue(context.Background(), "request_id", uuid.New().String())

	err = b.UploadSBOM(ctx, req)
	if err != nil {
		log.Fatal("UploadSBOM: ", err)
	}
}
