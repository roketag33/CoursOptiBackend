package dungeon

import (
	controller "dungeons/app/controllers/dungeon"
	service "dungeons/app/services/dungeon"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	dungeonService := service.New()
	dungeonController := controller.New(dungeonService)

	v1 := g.Group("/v1")
	{
		mj := v1.Group("/mj/dungeons")
		{
			mj.POST("", dungeonController.Create)
			mj.PUT("/:id", dungeonController.Update)
			mj.POST("/:id/publish", dungeonController.Publish)
			mj.POST("/:id/steps", dungeonController.AddStep)
			mj.PUT("/:id/steps/:stepId", dungeonController.UpdateStep)
			mj.PUT("/:id/steps/reorder", dungeonController.ReorderSteps)
		}

		dungeons := v1.Group("/dungeons")
		{
			dungeons.GET("", dungeonController.GetAll)
			dungeons.GET("/:id", dungeonController.GetByID)
		}
	}
}
