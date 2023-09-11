package app

import (
	"CRUDapp/internal/app/endpoint"
	"github.com/gin-gonic/gin"
	"log"
)

type App struct {
	endP   *endpoint.Endpoint
	router *gin.Engine
}

func New() (*App, error) {
	app := &App{}

	app.endP = endpoint.New()

	gin.SetMode(gin.ReleaseMode)
	app.router = gin.Default()
	err := app.router.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal(err)
	}

	app.router.GET("/entity/:id", app.endP.GetEntityByID)
	app.router.GET("/entities", app.endP.GetAllEntities)
	app.router.POST("/up", app.endP.PostEntity)
	app.router.PUT("/up/:id", app.endP.PutEntity)
	app.router.DELETE("/entity/:id", app.endP.DeleteEntityByID)

	return app, nil
}

func (app *App) Run() error {
	err := app.router.Run("0.0.0.0:8000")
	if err != nil {
		log.Fatal(err)
	}
	return err
}
