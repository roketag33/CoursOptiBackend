package run

import (
	"context"
	"dungeons/app/models"
	"dungeons/app/server"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func initTestEnv() *mongo.Database {
	_ = godotenv.Load("../../../.env")
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "mongodb://localhost:27017"
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dbHost).SetServerAPIOptions(serverAPI)
	client, _ := mongo.Connect(opts)
	db := client.Database("dungeons_test")
	
	srv := &server.Dungeons{
		Router:   gin.Default(),
		Database: db,
	}
	server.SetServer(srv)
	return db
}

func TestAttemptIdempotency(t *testing.T) {
	db := initTestEnv()
	db.Collection("run").Drop(context.TODO())
	db.Collection("bossstep").Drop(context.TODO())
	
	runCol := db.Collection("run")
	bossCol := db.Collection("bossstep")
	
	dungeonID := "test_dungeon"
	runID := "test_run"
	stepID := "test_step"
	
	boss := models.BossStep{
		CustomID:     stepID,
		DungeonID:    dungeonID,
		Order:        1,
		Lat:          0.0,
		Lon:          0.0,
		RadiusMeters: 500,
	}
	bossCol.InsertOne(context.TODO(), boss)
	
	run := models.Run{
		CustomID:    runID,
		DungeonID:   dungeonID,
		PlayerID:    "test_player",
		State:       "active",
		CurrentStep: 1,
		KilledSteps: []models.KilledStep{
			{BossStepID: stepID, KilledAt: time.Now()}, // Deja tué !
		},
	}
	runCol.InsertOne(context.TODO(), run)
	
	runService := New()
	
	// Attempt avec les bonnes coordonnees
	_, status, err := runService.Attempt(runID, stepID, 0.0, 0.0)
	
	if err == nil || status != 409 {
		t.Errorf("L'attempt aurait dû échouer avec status 409 (boss deja tué), got status = %d et err = %v", status, err)
	}
	
	if err != nil && err.Error() != "boss deja tué, pas de double récompense" {
		t.Errorf("Message d'erreur inattendu : %v", err)
	}
}

func TestAttemptProgressionOrder(t *testing.T) {
	db := initTestEnv()
	db.Collection("run").Drop(context.TODO())
	db.Collection("bossstep").Drop(context.TODO())
	
	runCol := db.Collection("run")
	bossCol := db.Collection("bossstep")
	
	dungeonID := "test_dungeon_2"
	runID := "test_run_2"
	stepID_1 := "test_step_1"
	stepID_2 := "test_step_2"
	
	bossCol.InsertOne(context.TODO(), models.BossStep{CustomID: stepID_1, DungeonID: dungeonID, Order: 1, Lat: 0.0, Lon: 0.0, RadiusMeters: 500})
	bossCol.InsertOne(context.TODO(), models.BossStep{CustomID: stepID_2, DungeonID: dungeonID, Order: 2, Lat: 0.0, Lon: 0.0, RadiusMeters: 500})
	
	// Le joueur est a l'etape 1
	runCol.InsertOne(context.TODO(), models.Run{
		CustomID:    runID,
		DungeonID:   dungeonID,
		State:       "active",
		CurrentStep: 1,
	})
	
	runService := New()
	
	// Essaie d'attaquer l'etape 2
	_, status, err := runService.Attempt(runID, stepID_2, 0.0, 0.0)
	if status != 409 || (err != nil && err.Error() != "WRONG_STEP_ORDER") {
		t.Errorf("L'attempt 2 aurait dû échouer pour Wrong Step Order. Got %d, %v", status, err)
	}
}
