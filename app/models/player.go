package models

import "time"

type PlayerID string

type Player struct {
	CustomID    string                 `bson:"customID"`
	DisplayName string                 `bson:"displayName" json:"display_name"`
	Password    string                 `bson:"password" json:"-"`
	GuildID     string                 `bson:"guildID,omitempty" json:"guild_id,omitempty"`
	Gold        int64                  `bson:"gold" json:"gold"`
	Stats       map[string]interface{} `bson:"stats" json:"stats"`
	Suspended   bool                   `bson:"suspended" json:"suspended"`
	CreatedAt   time.Time              `bson:"createdAt" json:"created_at"`
	UpdatedAt   time.Time              `bson:"updatedAt" json:"updated_at"`
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
