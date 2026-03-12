package models

import "time"

type ItemID string

type InventoryEntry struct {
	PlayerID  string    `bson:"playerID" json:"player_id"`
	ItemID    string    `bson:"itemID" json:"item_id"`
	Qty       int64     `bson:"qty" json:"qty"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updated_at"`
}

func (ie *InventoryEntry) Collection() string {
	return "inventory"
}

type ItemDef struct {
	CustomID    string    `bson:"customID" json:"id"`
	Type        string    `bson:"type" json:"type"`
	Rarity      string    `bson:"rarity" json:"rarity"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	Stats       map[string]interface{} `bson:"stats" json:"stats"`
	Tradable    bool      `bson:"tradable" json:"tradable"`
	BaseValue   int64     `bson:"baseValue" json:"base_value"`
	CreatedAt   time.Time `bson:"createdAt" json:"created_at"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"updated_at"`
}

func (i *ItemDef) Collection() string {
	return "item"
}

type InventoryResponse struct {
	PlayerID string             `json:"playerId"`
	Items    []InventoryItemDTO `json:"items"`
}

type InventoryItemDTO struct {
	ItemID string `json:"itemId"`
	Qty    int64  `json:"qty"`
}

type ItemDefResponse struct {
	ID          string         `json:"id"`
	Type        string         `json:"type"`
	Rarity      string         `json:"rarity"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Stats       map[string]any `json:"stats,omitempty"`
	Tradable    bool           `json:"tradable"`
	BaseValue   int64          `json:"baseValue,omitempty"`
}
