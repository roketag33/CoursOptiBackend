package run

import (
	"context"
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/server"
	"errors"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type AttemptRequest struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type AttemptResult struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Rewards *models.Reward  `json:"rewards,omitempty"`
	Gold    int64           `json:"gold_earned,omitempty"`
}

func (r *Run) Attempt(runID string, stepID string, lat float64, lon float64) (*AttemptResult, int, error) {
	srv := server.GetServer()

	run, err := r.GetByID(runID)
	if err != nil {
		return nil, 404, errors.New("run pas trouvée")
	}

	if run.State != "active" {
		return nil, 400, errors.New("cette run n'est plus active")
	}

	for _, killed := range run.KilledSteps {
		if killed.BossStepID == stepID {
			return nil, 409, errors.New("boss deja tué, pas de double récompense")
		}
	}

	bossCol := srv.Database.Collection("bossstep")
	var boss models.BossStep
	err = bossCol.FindOne(context.TODO(), bson.M{"customID": stepID, "dungeonID": run.DungeonID}).Decode(&boss)
	if err != nil {
		return nil, 404, errors.New("boss step pas trouvé")
	}

	if boss.Order != run.CurrentStep {
		return nil, 409, errors.New("WRONG_STEP_ORDER")
	}

	distance := functions.Haversine(lat, lon, boss.Lat, boss.Lon)
	if distance > boss.RadiusMeters {
		return nil, 409, errors.New("NOT_IN_RANGE")
	}

	won := rand.Intn(2) == 1
	if !won {
		return &AttemptResult{
			Success: false,
			Message: "Tu as perdu le combat contre " + boss.Name + "... Réessaie !",
		}, 200, nil
	}

	var totalGold int64
	for _, reward := range boss.Rewards {
		totalGold += reward.Gold
	}

	runCol := srv.Database.Collection("run")
	playerCol := srv.Database.Collection("player")
	inventoryCol := srv.Database.Collection("inventory")

	killedStep := models.KilledStep{
		BossStepID: stepID,
		KilledAt:   time.Now(),
	}

	dungeonCol := srv.Database.Collection("dungeon")
	var dg models.Dungeon
	dungeonCol.FindOne(context.TODO(), bson.M{"customID": run.DungeonID}).Decode(&dg)

	newStep := run.CurrentStep + 1
	newState := run.State
	if newStep > len(dg.Bosses) {
		newState = "completed"
	}

	updateRun := bson.M{
		"$push": bson.M{"killedSteps": killedStep},
		"$set":  bson.M{"currentStep": newStep, "state": newState},
	}
	if newState == "completed" {
		now := time.Now()
		updateRun["$set"] = bson.M{"currentStep": newStep, "state": newState, "endedAt": now}
	}
	_, err = runCol.UpdateOne(context.TODO(), bson.M{"customID": runID}, updateRun)
	if err != nil {
		return nil, 500, err
	}

	if totalGold > 0 {
		playerCol.UpdateOne(context.TODO(), bson.M{"customID": run.PlayerID}, bson.M{
			"$inc": bson.M{"gold": totalGold},
		})
	}

	for _, reward := range boss.Rewards {
		if reward.ItemID == "" || reward.Qty <= 0 {
			continue
		}

		var existing models.InventoryEntry
		err = inventoryCol.FindOne(context.TODO(), bson.M{
			"playerID": run.PlayerID,
			"itemID":   reward.ItemID,
		}).Decode(&existing)

		if err != nil {
			inventoryCol.InsertOne(context.TODO(), models.InventoryEntry{
				PlayerID:  run.PlayerID,
				ItemID:    reward.ItemID,
				Qty:       int64(reward.Qty),
				UpdatedAt: time.Now(),
			})
		} else {
			inventoryCol.UpdateOne(context.TODO(), bson.M{
				"playerID": run.PlayerID,
				"itemID":   reward.ItemID,
			}, bson.M{
				"$inc": bson.M{"qty": reward.Qty},
				"$set": bson.M{"updatedAt": time.Now()},
			})
		}
	}

	return &AttemptResult{
		Success: true,
		Message: "Tu as vaincu " + boss.Name + " !",
		Gold:    totalGold,
	}, 200, nil
}
