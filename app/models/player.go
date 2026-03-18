package models

import "time"

type PlayerID string

type Player struct {
	CustomID    string                 `bson:"customID"`
	DisplayName string                 `bson:"displayName"`
	Password    string                 `bson:"password" json:"-"`
	Gold        int64                  `bson:"gold"`
	Stats       map[string]interface{} `bson:"stats"`
	CreatedAt   time.Time              `bson:"createdAt"`
	UpdatedAt   time.Time              `bson:"updatedAt"`
}

type PlayerResponse struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"displayName"`
	Wallet      Wallet    `json:"wallet"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Wallet struct {
	Gold int64 `json:"gold"` // int64 pour éviter les soucis quand ça grossit
}

// Collection Mongodb collection
func (p *Player) Collection() string {
	return "player"
}
