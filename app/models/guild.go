package models

import "time"

type Guild struct {
	CustomID    string    `bson:"customID" json:"id"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	CreatorID   string    `bson:"creatorID" json:"creator_id"`
	CreatedAt   time.Time `bson:"createdAt" json:"created_at"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"updated_at"`
}

func (g *Guild) Collection() string {
	return "guild"
}
