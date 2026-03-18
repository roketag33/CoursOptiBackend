package auth

import (
	controller "dungeons/app/controllers/auth"
	"dungeons/app/middlewares"
	service "dungeons/app/services/auth"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	authService := service.New()
	authController := controller.New(authService)

	v1 := g.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}
		
		me := v1.Group("/me")
		me.Use(middlewares.AuthMiddleware())
		{
			me.GET("", authController.Me)
		}
	}
}
