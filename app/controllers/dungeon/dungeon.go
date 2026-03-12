package dungeon

import (
	"dungeons/app/controllers/common"
	"dungeons/app/models"
	service "dungeons/app/services/dungeon"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Dungeon struct {
	DungeonService *service.Dungeon
}

func New(s *service.Dungeon) *Dungeon {
	return &Dungeon{DungeonService: s}
}

func (d *Dungeon) Create(ctx *gin.Context) {
	var in models.Dungeon
	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "dungeon.Create.BadRequest", err))
		return
	}

	dungeon, err := d.DungeonService.Create(&in)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, "dungeon.Create.Error", err))
		return
	}

	common.SendResponse(ctx, http.StatusCreated, dungeon)
}

func (d *Dungeon) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	dungeon, err := d.DungeonService.GetByID(id)
	if err != nil {
		common.SendResponse(ctx, http.StatusNotFound, models.KnownError(http.StatusNotFound, "dungeon.Get.NotFound", err))
		return
	}

	common.SendResponse(ctx, http.StatusOK, dungeon)
}

func (d *Dungeon) GetAll(ctx *gin.Context) {
	dungeons, err := d.DungeonService.GetAll()
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, "dungeon.GetAll.Error", err))
		return
	}

	common.SendResponse(ctx, http.StatusOK, dungeons)
}

func (d *Dungeon) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var in models.Dungeon
	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "dungeon.Update.BadRequest", err))
		return
	}

	err := d.DungeonService.Update(id, &in)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "dungeon.Update.Error", err))
		return
	}

	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, "dungeon.Update.OK", "dungeon mis a jour"))
}

func (d *Dungeon) Publish(ctx *gin.Context) {
	id := ctx.Param("id")

	err := d.DungeonService.Publish(id)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "dungeon.Publish.Error", err))
		return
	}

	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, "dungeon.Publish.OK", "dungeon publié"))
}

func (d *Dungeon) AddStep(ctx *gin.Context) {
	dungeonID := ctx.Param("id")
	var in models.BossStep
	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "dungeon.AddStep.BadRequest", err))
		return
	}

	step, err := d.DungeonService.AddStep(dungeonID, &in)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, "dungeon.AddStep.Error", err))
		return
	}

	common.SendResponse(ctx, http.StatusCreated, step)
}

func (d *Dungeon) UpdateStep(ctx *gin.Context) {
	dungeonID := ctx.Param("id")
	stepID := ctx.Param("stepId")
	var in models.BossStep
	if err := ctx.BindJSON(&in); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "dungeon.UpdateStep.BadRequest", err))
		return
	}

	err := d.DungeonService.UpdateStep(dungeonID, stepID, &in)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "dungeon.UpdateStep.Error", err))
		return
	}

	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, "dungeon.UpdateStep.OK", "step mis a jour"))
}

type ReorderRequest struct {
	StepIDs []string `json:"step_ids"`
}

func (d *Dungeon) ReorderSteps(ctx *gin.Context) {
	dungeonID := ctx.Param("id")
	var req ReorderRequest
	if err := ctx.BindJSON(&req); err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "dungeon.Reorder.BadRequest", err))
		return
	}

	err := d.DungeonService.ReorderSteps(dungeonID, req.StepIDs)
	if err != nil {
		common.SendResponse(ctx, http.StatusBadRequest, models.KnownError(http.StatusBadRequest, "dungeon.Reorder.Error", err))
		return
	}

	common.SendResponse(ctx, http.StatusOK, models.Success(http.StatusOK, "dungeon.Reorder.OK", "steps réordonnés"))
}
