package main

import (
	"context"
	"dungeons/app/models"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("pas de .env trouvé, on utilise les variables d'environnement")
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		fmt.Println("DB_HOST est vide !!")
		os.Exit(1)
	}

	fmt.Println("Connexion a MongoDB...")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dbHost).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)
	if err != nil {
		fmt.Println("Erreur connexion:", err)
		os.Exit(1)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database("dungeons")
	fmt.Println("Connecté a la base dungeons")

	fmt.Println("On drop les collections...")
	collections := []string{"item", "player", "dungeon", "bossstep", "run", "inventory"}
	for _, col := range collections {
		db.Collection(col).Drop(context.TODO())
	}
	fmt.Println("Collections supprimées")

	seedItems(db)
	seedPlayers(db)
	seedDungeon(db)
	seedInventory(db)

	fmt.Println("")
	fmt.Println("===== Seed terminé ! =====")
}

func seedItems(db *mongo.Database) {
	fmt.Println("Insertion des items...")
	col := db.Collection("item")

	items := []interface{}{
		models.ItemDef{
			CustomID:    "item-epee-rouille",
			Type:        "weapon",
			Rarity:      "common",
			Name:        "Epee Rouillee",
			Description: "Une vieille epee qui a connu des jours meilleurs",
			Stats:       map[string]interface{}{"attack": 5, "durability": 20},
			Tradable:    true,
			BaseValue:   10,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		models.ItemDef{
			CustomID:    "item-potion-soin",
			Type:        "consumable",
			Rarity:      "common",
			Name:        "Potion de Soin",
			Description: "Restaure un peu de vie",
			Stats:       map[string]interface{}{"heal": 25},
			Tradable:    true,
			BaseValue:   5,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		models.ItemDef{
			CustomID:    "item-bouclier-bois",
			Type:        "weapon",
			Rarity:      "uncommon",
			Name:        "Bouclier en Bois",
			Description: "Un bouclier basique en bois",
			Stats:       map[string]interface{}{"defense": 8, "durability": 30},
			Tradable:    true,
			BaseValue:   15,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	_, err := col.InsertMany(context.TODO(), items)
	if err != nil {
		fmt.Println("Erreur insertion items:", err)
		return
	}
	fmt.Println("3 items insérés")
}

func seedPlayers(db *mongo.Database) {
	fmt.Println("Insertion des players...")
	col := db.Collection("player")

	players := []interface{}{
		models.Player{
			CustomID:    "player-alexandre",
			DisplayName: "Alexandre",
			Gold:        100,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		models.Player{
			CustomID:    "player-marie",
			DisplayName: "Marie",
			Gold:        150,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	_, err := col.InsertMany(context.TODO(), players)
	if err != nil {
		fmt.Println("Erreur insertion players:", err)
		return
	}
	fmt.Println("2 players insérés")
}

func seedDungeon(db *mongo.Database) {
	fmt.Println("Insertion du dungeon et des boss steps...")

	dungeonCol := db.Collection("dungeon")
	bossCol := db.Collection("bossstep")

	dungeon := models.Dungeon{
		CustomID:    "dungeon-paris-catacombes",
		Title:       "Les Catacombes de Paris",
		Description: "Explorez les souterrains de Paris et affrontez les gardiens des ombres",
		CreatedBy:   "mj-001",
		Area:        models.Area{Name: "Paris Centre"},
		Bosses:      []string{"boss-squelette", "boss-spectre", "boss-liche"},
		Status:      "published",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := dungeonCol.InsertOne(context.TODO(), dungeon)
	if err != nil {
		fmt.Println("Erreur insertion dungeon:", err)
		return
	}
	fmt.Println("1 dungeon inséré")

	bosses := []interface{}{
		models.BossStep{
			CustomID:        "boss-squelette",
			DungeonID:       "dungeon-paris-catacombes",
			Order:           1,
			Name:            "Squelette Ancien",
			Lat:             48.8340,
			Lon:             2.3321,
			RadiusMeters:    50,
			ZoneDescription: "Pres de la sortie des Catacombes, Avenue du Colonel Henri Rol-Tanguy",
			Difficulty:      3,
			Rewards: []models.Reward{
				{ItemID: "item-epee-rouille", Qty: 1, Gold: 20},
			},
		},
		models.BossStep{
			CustomID:        "boss-spectre",
			DungeonID:       "dungeon-paris-catacombes",
			Order:           2,
			Name:            "Spectre du Metro",
			Lat:             48.8332,
			Lon:             2.3317,
			RadiusMeters:    30,
			ZoneDescription: "Station Denfert-Rochereau, dans le couloir sombre",
			Difficulty:      5,
			Rewards: []models.Reward{
				{ItemID: "item-potion-soin", Qty: 2, Gold: 35},
			},
		},
		models.BossStep{
			CustomID:        "boss-liche",
			DungeonID:       "dungeon-paris-catacombes",
			Order:           3,
			Name:            "La Liche des Profondeurs",
			Lat:             48.8362,
			Lon:             2.3350,
			RadiusMeters:    40,
			ZoneDescription: "Place Denfert-Rochereau, pres du Lion de Belfort",
			Difficulty:      8,
			Rewards: []models.Reward{
				{ItemID: "item-bouclier-bois", Qty: 1, Gold: 75},
				{ItemID: "item-potion-soin", Qty: 3, Gold: 0},
			},
		},
	}

	_, err = bossCol.InsertMany(context.TODO(), bosses)
	if err != nil {
		fmt.Println("Erreur insertion boss steps:", err)
		return
	}
	fmt.Println("3 boss steps insérés")
}

func seedInventory(db *mongo.Database) {
	fmt.Println("Insertion de l'inventaire...")
	col := db.Collection("inventory")

	entries := []interface{}{
		models.InventoryEntry{
			PlayerID:  "player-alexandre",
			ItemID:    "item-potion-soin",
			Qty:       3,
			UpdatedAt: time.Now(),
		},
		models.InventoryEntry{
			PlayerID:  "player-marie",
			ItemID:    "item-epee-rouille",
			Qty:       1,
			UpdatedAt: time.Now(),
		},
	}

	_, err := col.InsertMany(context.TODO(), entries)
	if err != nil {
		fmt.Println("Erreur insertion inventaire:", err)
		return
	}
	fmt.Println("2 entrées d'inventaire insérées")
}
