package auction

import (
	controller "dungeons/app/controllers/auction"
	service "dungeons/app/services/auction"

	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine) {
	auctionService := service.New()
	auctionController := controller.New(auctionService)

	v1 := g.Group("/v1")
	{
		act := v1.Group("/auction")
		{
			act.GET("/listings", auctionController.GetActiveListings)
			act.POST("/listings", auctionController.CreateListing)
			act.POST("/listings/:id/buy", auctionController.BuyListing)
			act.POST("/listings/:id/cancel", auctionController.CancelListing)
		}
	}
}
