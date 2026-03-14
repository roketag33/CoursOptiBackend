package inventory

import (
	"dungeons/app/controllers/common"
	"dungeons/app/models"
	service "dungeons/app/services/inventory"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Inventory struct {
	InventoryService *service.Inventory
}

func New(s *service.Inventory) *Inventory {
	return &Inventory{InventoryService: s}
}

func (i *Inventory) Get(ctx *gin.Context) {
	playerID := ctx.GetHeader("X-Player-Id")
	if playerID == "" {
		common.SendResponse(ctx, http.StatusBadRequest, models.Success(http.StatusBadRequest, "inventory.Get.BadRequest", "Header X-Player-Id requis"))
		return
	}

	entries, err := i.InventoryService.GetByPlayerID(playerID)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, "inventory.Get.Error", err))
		return
	}

	if entries == nil {
		entries = []models.InventoryEntry{}
	}

	common.SendResponse(ctx, http.StatusOK, entries)
}
