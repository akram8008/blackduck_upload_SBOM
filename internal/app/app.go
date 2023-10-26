package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/akram8008/blackduck_upload_SBOM/internal/config"
	"github.com/akram8008/blackduck_upload_SBOM/pkg/blackduck"
	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/net/context"
)

type App struct {
	blackduck *blackduck.BlackDuckComponents
	l         *log.Logger
}

func New(config *config.Config, l *log.Logger) *App {
	return &App{
		blackduck: blackduck.NewBlackDuckComponents(
			l,
			config.BlackduckUrl, config.BlackduckToken),
		l: l,
	}
}

func (app *App) Run() {

	router := gin.Default()

	router.GET("/upload/sbom/:projectName/:versionName", app.uploadSbom)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/ping", app.Ping())

	log.Fatalln(router.Run("localhost:8080"))
}

// @Summary upload sbom
// @Description upload sbom
// @Router /upload/sbom/{projectName}/{versionName} [GET]
// @Tags SBOM
// @Accept json
// @Produce json
// @Param projectName path string true "project name"
// @Param versionName path string true "project version"
// @Success 200 {object} view
// @Failure 400 {object} view
// @Failure 500 {object} view
func (app *App) uploadSbom(c *gin.Context) {
	projectName := c.Param("projectName")
	versionName := c.Param("versionName")

	if projectName == "" || versionName == "" {
		c.JSON(http.StatusBadRequest, view{
			Status: "failed",
			Msg:    "projectName and versionName must be provided",
		})
		return
	}

	f, err := os.Open("./file/bom.xml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, view{
			Status: "failed",
			Msg:    err.Error(),
		})
		return

	}
	defer f.Close()

	ctx := context.WithValue(context.Background(), "request_id", uuid.New().String())

	req := blackduck.RequestUploadSBOM{
		File:           f,
		FileName:       "bom.xml",
		ProjectName:    projectName,
		ProjectVersion: versionName,
	}

	err = app.blackduck.UploadSBOM(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, view{
			Status: "failed",
			Msg:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, view{
		Status: "success",
		Msg: fmt.Sprintf("sbom file for project: %s, version: %s successfully uploaded",
			projectName, versionName),
	})
}

// @Summary ping
// @Description ping
// @Router /ping [GET]
// @Tags SBOM
// @Accept json
// @Produce json
// @Success 200 {object} view
// @Failure 400 {object} view
// @Failure 500 {object} view
func (app *App) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, view{
			Status: "success",
			Msg:    "pong",
		})
	}
}
