package inventory

import (
	"context"
	"dungeons/app/models"
	"dungeons/app/server"

	"go.mongodb.org/mongo-driver/bson"
)

type Inventory struct{}

func New() *Inventory {
	return &Inventory{}
}

func (i *Inventory) GetByPlayerID(playerID string) ([]models.InventoryEntry, error) {
	srv := server.GetServer()
	col := srv.Database.Collection("inventory")

	var entries []models.InventoryEntry
	cursor, err := col.Find(context.TODO(), bson.M{"playerID": playerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var entry models.InventoryEntry
		if err := cursor.Decode(&entry); err != nil {
			return nil, err
		}
		// On ignore les entrees avec 0 en quantite (si un item a ete utilise ou vendu completement)
		if entry.Qty > 0 {
			entries = append(entries, entry)
		}
	}

	return entries, nil
}
