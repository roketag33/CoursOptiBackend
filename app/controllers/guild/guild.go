package guild

import (
	"dungeons/app/controllers/common"
	"dungeons/app/models"
	service "dungeons/app/services/guild"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Guild struct {
	GuildService *service.Guild
}

func New(s *service.Guild) *Guild {
	return &Guild{GuildService: s}
}

type CreateGuildRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (g *Guild) Create(ctx *gin.Context) {
	playerID, exists := ctx.Get("playerID")
	if !exists {
		common.SendResponse(ctx, http.StatusUnauthorized, models.Success(http.StatusUnauthorized, "guild.Create.Unauthorized", "Pas connecté"))
		return
	}

	var req CreateGuildRequest
	if err := ctx.BindJSON(&req); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "guild.Create.BadRequest", err))
		return
	}

	if req.Name == "" {
		common.SendResponse(ctx, http.StatusBadRequest, models.Success(http.StatusBadRequest, "guild.Create.Invalid", "Le nom de guilde est requis"))
		return
	}

	guild, err := g.GuildService.Create(req.Name, req.Description, playerID.(string))
	if err != nil {
		common.SendResponse(ctx, http.StatusConflict, models.KnownError(http.StatusConflict, "guild.Create.Conflict", err))
		return
	}

	common.SendResponse(ctx, http.StatusCreated, guild)
}

func (g *Guild) Join(ctx *gin.Context) {
	playerID, exists := ctx.Get("playerID")
	if !exists {
		common.SendResponse(ctx, http.StatusUnauthorized, models.Success(http.StatusUnauthorized, "guild.Join.Unauthorized", "Pas connecté"))
		return
	}

	guildID := ctx.Param("id")
	err := g.GuildService.Join(guildID, playerID.(string))
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "guild.Join.Error", err))
		return
	}

	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, "guild.Join.Success", "Tu as rejoint la guilde"))
}
