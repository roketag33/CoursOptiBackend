package models

type Reward struct {
	ItemID string `bson:"itemID" json:"item_id"`
	Qty    int    `bson:"qty" json:"qty"`
	Gold   int64  `bson:"gold" json:"gold"`
}

type BossStep struct {
	CustomID        string   `bson:"customID" json:"id"`
	DungeonID       string   `bson:"dungeonID" json:"dungeon_id"`
	Order           int      `bson:"order" json:"order"`
	Name            string   `bson:"name" json:"name"`
	Lat             float64  `bson:"lat" json:"lat"`
	Lon             float64  `bson:"lon" json:"lon"`
	RadiusMeters    float64  `bson:"radiusMeters" json:"radius_meters"`
	ZoneDescription string   `bson:"zoneDescription" json:"zone_description"`
	Difficulty      int      `bson:"difficulty" json:"difficulty"`
	Rewards         []Reward `bson:"rewards" json:"rewards"`
}

func (b *BossStep) Collection() string {
	return "bossstep"
}
