package run

import (
	"dungeons/app/controllers/common"
	"dungeons/app/models"
	service "dungeons/app/services/run"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Run struct {
	RunService *service.Run
}

func New(s *service.Run) *Run {
	return &Run{RunService: s}
}

type CreateRunRequest struct {
	DungeonID string `json:"dungeon_id"`
	PlayerID  string `json:"player_id"`
}

func (r *Run) Create(ctx *gin.Context) {
	var req CreateRunRequest
	if err := ctx.BindJSON(&req); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "run.Create.BadRequest", err))
		return
	}

	if req.DungeonID == "" || req.PlayerID == "" {
		common.SendResponse(ctx, http.StatusBadRequest, models.Success(http.StatusBadRequest, "run.Create.BadRequest", "dungeon_id et player_id requis"))
		return
	}

	run, err := r.RunService.Create(req.DungeonID, req.PlayerID)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "run.Create.Error", err))
		return
	}

	common.SendResponse(ctx, http.StatusCreated, run)
}

func (r *Run) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	run, err := r.RunService.GetByID(id)
	if err != nil {
		common.SendResponse(ctx, http.StatusNotFound, models.KnownError(http.StatusNotFound, "run.Get.NotFound", err))
		return
	}

	common.SendResponse(ctx, http.StatusOK, run)
}

func (r *Run) GetAll(ctx *gin.Context) {
	playerID := ctx.Query("player_id")

	var runs []models.Run
	var err error

	if playerID != "" {
		runs, err = r.RunService.GetByPlayerID(playerID)
	} else {
		runs, err = r.RunService.GetAll()
	}

	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, "run.GetAll.Error", err))
		return
	}

	if runs == nil {
		runs = []models.Run{}
	}

	common.SendResponse(ctx, http.StatusOK, runs)
}

func (r *Run) Attempt(ctx *gin.Context) {
	runID := ctx.Param("id")
	stepID := ctx.Param("stepId")

	var req service.AttemptRequest
	if err := ctx.BindJSON(&req); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "run.Attempt.BadRequest", err))
		return
	}

	result, status, err := r.RunService.Attempt(runID, stepID, req.Lat, req.Lon)
	if err != nil {
		if status == 404 {
			common.SendResponse(ctx, http.StatusNotFound, models.KnownError(http.StatusNotFound, "run.Attempt.NotFound", err))
		} else if status == 409 {
			common.SendResponse(ctx, http.StatusConflict, models.KnownError(http.StatusConflict, "run.Attempt.Conflict", err))
		} else {
			common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, "run.Attempt.Error", err))
		}
		return
	}

	common.SendResponse(ctx, http.StatusOK, result)
}
