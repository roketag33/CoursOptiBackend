package auction

import (
	"dungeons/app/controllers/common"
	"dungeons/app/models"
	service "dungeons/app/services/auction"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Auction struct {
	AuctionService *service.Auction
}

func New(s *service.Auction) *Auction {
	return &Auction{AuctionService: s}
}

type CreateListingReq struct {
	ItemID       string `json:"item_id"`
	Qty          int    `json:"qty"`
	PricePerUnit int64  `json:"price_per_unit"`
}

func (a *Auction) CreateListing(ctx *gin.Context) {
	playerID := ctx.GetHeader("X-Player-Id")
	if playerID == "" {
		common.SendResponse(ctx, http.StatusBadRequest, models.Success(http.StatusBadRequest, "auction.Create.BadRequest", "Header X-Player-Id requis"))
		return
	}

	var req CreateListingReq
	if err := ctx.BindJSON(&req); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "auction.Create.BadRequest", err))
		return
	}

	if req.Qty <= 0 || req.PricePerUnit < 0 {
		common.SendResponse(ctx, http.StatusBadRequest, models.Success(http.StatusBadRequest, "auction.Create.BadRequest", "qty ou prix invalides"))
		return
	}

	listing, err := a.AuctionService.CreateListing(playerID, req.ItemID, req.Qty, req.PricePerUnit)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, "auction.Create.Error", err))
		return
	}

	common.SendResponse(ctx, http.StatusCreated, listing)
}

func (a *Auction) GetActiveListings(ctx *gin.Context) {
	listings, err := a.AuctionService.GetActiveListings()
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, "auction.Get.Error", err))
		return
	}

	if listings == nil {
		listings = []models.Listing{}
	}

	common.SendResponse(ctx, http.StatusOK, listings)
}

func (a *Auction) BuyListing(ctx *gin.Context) {
	playerID := ctx.GetHeader("X-Player-Id")
	if playerID == "" {
		common.SendResponse(ctx, http.StatusBadRequest, models.Success(http.StatusBadRequest, "auction.Buy.BadRequest", "Header X-Player-Id requis"))
		return
	}

	id := ctx.Param("id")
	err := a.AuctionService.BuyListing(playerID, id)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, "auction.Buy.Error", err))
		return
	}

	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, "auction.Buy.OK", "Achat reussi"))
}

func (a *Auction) CancelListing(ctx *gin.Context) {
	playerID := ctx.GetHeader("X-Player-Id")
	if playerID == "" {
		common.SendResponse(ctx, http.StatusBadRequest, models.Success(http.StatusBadRequest, "auction.Cancel.BadRequest", "Header X-Player-Id requis"))
		return
	}

	id := ctx.Param("id")
	err := a.AuctionService.CancelListing(playerID, id)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, "auction.Cancel.Error", err))
		return
	}

	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, "auction.Cancel.OK", "Listing annule"))
}
