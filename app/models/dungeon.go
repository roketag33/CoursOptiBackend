package models

import "time"

type Area struct {
	Name string `bson:"name" json:"name"`
}

type Dungeon struct {
	CustomID    string    `bson:"customID" json:"id"`
	Title       string    `bson:"title" json:"title"`
	Description string    `bson:"description" json:"description"`
	CreatedBy   string    `bson:"createdBy" json:"created_by"`
	Area        Area      `bson:"area" json:"area"`
	Bosses      []string  `bson:"bosses" json:"bosses"`
	Status      string    `bson:"status" json:"status"`
	CreatedAt   time.Time `bson:"createdAt" json:"created_at"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"updated_at"`
}

func (d *Dungeon) Collection() string {
	return "dungeon"
}
