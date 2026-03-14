package models

import "time"

type Listing struct {
	CustomID     string    `bson:"customID" json:"id"`
	SellerID     string    `bson:"sellerID" json:"seller_id"`
	ItemID       string    `bson:"itemID" json:"item_id"`
	Qty          int       `bson:"qty" json:"qty"`
	PricePerUnit int64     `bson:"pricePerUnit" json:"price_per_unit"`
	Status       string    `bson:"status" json:"status"`
	CreatedAt    time.Time `bson:"createdAt" json:"created_at"`
	ExpiresAt    time.Time `bson:"expiresAt" json:"expires_at"`
}

func (l *Listing) Collection() string {
	return "listing"
}
