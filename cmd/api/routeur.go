package main

import (
	controllers "dungeons/app/controllers/common"
	"dungeons/app/middlewares"
	routes "dungeons/app/routes/common"

	"github.com/gin-gonic/gin"
)

// init the router
func setupRouter() *gin.Engine {
	router := routes.SetupRouter()
	router.Use(middlewares.RateLimitMiddleware())
	router.GET("/ping", controllers.Ping)
	router.GET("/version", controllers.Version)

	return router
}
