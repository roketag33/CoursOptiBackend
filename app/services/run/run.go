package run

import (
	"context"
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/server"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Run struct{}

func New() *Run {
	return &Run{}
}

func (r *Run) Create(dungeonID string, playerID string) (*models.Run, error) {
	srv := server.GetServer()
	col := srv.Database.Collection("run")
	dungeonCol := srv.Database.Collection("dungeon")

	var dungeon models.Dungeon
	err := dungeonCol.FindOne(context.TODO(), bson.M{"customID": dungeonID}).Decode(&dungeon)
	if err != nil {
		return nil, errors.New("dungeon pas trouvé")
	}

	if dungeon.Status != "published" {
		return nil, errors.New("le dungeon n'est pas publié")
	}

	var existingRun models.Run
	err = col.FindOne(context.TODO(), bson.M{
		"dungeonID": dungeonID,
		"playerID":  playerID,
		"state":     "active",
	}).Decode(&existingRun)
	if err == nil {
		return nil, errors.New("tu as deja une run active sur ce dungeon")
	}

	newRun := models.Run{
		CustomID:    functions.NewUUID(),
		DungeonID:   dungeonID,
		PlayerID:    playerID,
		State:       "active",
		CurrentStep: 1,
		KilledSteps: []models.KilledStep{},
		StartedAt:   time.Now(),
	}

	_, err = col.InsertOne(context.TODO(), newRun)
	if err != nil {
		return nil, err
	}

	return &newRun, nil
}

func (r *Run) GetByID(id string) (*models.Run, error) {
	srv := server.GetServer()
	col := srv.Database.Collection("run")

	var run models.Run
	err := col.FindOne(context.TODO(), bson.M{"customID": id}).Decode(&run)
	if err != nil {
		return nil, errors.New("run pas trouvée")
	}

	return &run, nil
}

func (r *Run) GetByPlayerID(playerID string) ([]models.Run, error) {
	srv := server.GetServer()
	col := srv.Database.Collection("run")

	var runs []models.Run
	cursor, err := col.Find(context.TODO(), bson.M{"playerID": playerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var run models.Run
		if err := cursor.Decode(&run); err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}

	return runs, nil
}

func (r *Run) GetAll() ([]models.Run, error) {
	srv := server.GetServer()
	col := srv.Database.Collection("run")

	var runs []models.Run
	cursor, err := col.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var run models.Run
		if err := cursor.Decode(&run); err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}

	return runs, nil
}
