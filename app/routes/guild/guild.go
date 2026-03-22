package guild

import (
	controller "dungeons/app/controllers/guild"
	"dungeons/app/middlewares"
	service "dungeons/app/services/guild"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	guildService := service.New()
	guildController := controller.New(guildService)

	v1 := g.Group("/v1")
	{
		guilds := v1.Group("/guilds")
		guilds.Use(middlewares.AuthMiddleware())
		{
			guilds.POST("", guildController.Create)
			guilds.POST("/:id/join", guildController.Join)
		}
	}
}
