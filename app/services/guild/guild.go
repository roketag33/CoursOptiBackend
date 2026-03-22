package guild

import (
	"context"
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/server"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Guild struct{}

func New() *Guild {
	return &Guild{}
}

func (g *Guild) Create(name, description, creatorID string) (*models.Guild, error) {
	srv := server.GetServer()
	guildCol := srv.Database.Collection("guild")
	playerCol := srv.Database.Collection("player")

	var existing models.Guild
	err := guildCol.FindOne(context.TODO(), bson.M{"name": name}).Decode(&existing)
	if err == nil {
		return nil, errors.New("ce nom de guilde est deja pris")
	}

	guild := models.Guild{
		CustomID:    functions.NewUUID(),
		Name:        name,
		Description: description,
		CreatorID:   creatorID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = guildCol.InsertOne(context.TODO(), guild)
	if err != nil {
		return nil, err
	}

	_, err = playerCol.UpdateOne(context.TODO(), bson.M{"customID": creatorID}, bson.M{
		"$set": bson.M{"guildID": guild.CustomID},
	})

	if err != nil {
		return nil, err
	}

	return &guild, nil
}

func (g *Guild) Join(guildID, playerID string) error {
	srv := server.GetServer()
	guildCol := srv.Database.Collection("guild")
	playerCol := srv.Database.Collection("player")

	var guild models.Guild
	err := guildCol.FindOne(context.TODO(), bson.M{"customID": guildID}).Decode(&guild)
	if err != nil {
		return errors.New("guilde introuvable")
	}

	var player models.Player
	err = playerCol.FindOne(context.TODO(), bson.M{"customID": playerID}).Decode(&player)
	if err != nil {
		return errors.New("joueur introuvable")
	}

	if player.GuildID != "" {
		return errors.New("tu es deja dans une guilde")
	}

	_, err = playerCol.UpdateOne(context.TODO(), bson.M{"customID": playerID}, bson.M{
		"$set": bson.M{"guildID": guild.CustomID},
	})

	return err
}
