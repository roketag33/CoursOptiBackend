package auth

import (
	"dungeons/app/controllers/common"
	"dungeons/app/models"
	service "dungeons/app/services/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	AuthService *service.Auth
}

func New(s *service.Auth) *Auth {
	return &Auth{AuthService: s}
}

type AuthRequest struct {
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}

func (a *Auth) Register(ctx *gin.Context) {
	var req AuthRequest
	if err := ctx.BindJSON(&req); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "auth.Register.BadRequest", err))
		return
	}

	if len(req.Password) < 4 || req.DisplayName == "" {
		common.SendResponse(ctx, http.StatusBadRequest, models.Success(http.StatusBadRequest, "auth.Register.Invalid", "Nom d'utilisateur et mdp(min 4) requis"))
		return
	}

	player, err := a.AuthService.Register(req.DisplayName, req.Password)
	if err != nil {
		common.SendResponse(ctx, http.StatusConflict, models.KnownError(http.StatusConflict, "auth.Register.Conflict", err))
		return
	}

	common.SendResponse(ctx, http.StatusCreated, player)
}

func (a *Auth) Login(ctx *gin.Context) {
	var req AuthRequest
	if err := ctx.BindJSON(&req); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "auth.Login.BadRequest", err))
		return
	}

	token, err := a.AuthService.Login(req.DisplayName, req.Password)
	if err != nil {
		common.SendResponse(ctx, http.StatusUnauthorized, models.KnownError(http.StatusUnauthorized, "auth.Login.Unauthorized", err))
		return
	}

	common.SendResponse(ctx, http.StatusOK, gin.H{"access_token": token})
}

func (a *Auth) Me(ctx *gin.Context) {
	playerID, exists := ctx.Get("playerID")
	if !exists {
		common.SendResponse(ctx, http.StatusUnauthorized, models.Success(http.StatusUnauthorized, "auth.Me.Unauthorized", "Pas connecté"))
		return
	}

	common.SendResponse(ctx, http.StatusOK, gin.H{"customID": playerID})
}
