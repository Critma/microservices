package handlers

import (
	"net/http"

	"github.com/critma/prodfiles/cmd/api/config"
	"github.com/gin-gonic/gin"
)

func AddHandlers(router *gin.Engine, app config.Application) {
	router.GET("/status", func(ctx *gin.Context) {
		message := "Welcome to api"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	filesHandler := NewFiles(app.Store, app.Logger, app.Config)

	api := router.Group("/api")
	v1 := api.Group("/v1")
	{
		images := v1.Group("/images")
		{
			// images.StaticFS("/", http.Dir(app.Config.BasePath))
			images.GET("/:id/:filename", filesHandler.GetFile)
			images.POST("/:id/:filename", filesHandler.UploadREST)
			images.POST("/", filesHandler.UploadMultipart)
		}
	}
}
