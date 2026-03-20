package auction

import (
	"context"
	"dungeons/app/models"
	"dungeons/app/server"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func initTestEnv() *mongo.Database {
	_ = godotenv.Load("../../../.env")
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "mongodb://localhost:27017"
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dbHost).SetServerAPIOptions(serverAPI)
	client, _ := mongo.Connect(opts)
	db := client.Database("dungeons_test")
	
	srv := &server.Dungeons{
		Router:   gin.Default(),
		Database: db,
	}
	server.SetServer(srv)
	return db
}

func TestAuctionTransactions(t *testing.T) {
	db := initTestEnv()
	
	// Nettoyer db
	db.Collection("player").Drop(context.TODO())
	db.Collection("inventory").Drop(context.TODO())
	db.Collection("listing").Drop(context.TODO())
	
	playerCol := db.Collection("player")
	invCol := db.Collection("inventory")
	
	sellerID := "test_seller"
	buyerID := "test_buyer"
	itemID := "test_sword"
	
	playerCol.InsertOne(context.TODO(), models.Player{CustomID: sellerID, Gold: 0})
	playerCol.InsertOne(context.TODO(), models.Player{CustomID: buyerID, Gold: 100})
	
	invCol.InsertOne(context.TODO(), models.InventoryEntry{PlayerID: sellerID, ItemID: itemID, Qty: 5, UpdatedAt: time.Now()})
	
	act := New()
	
	// Etape 1 : Créer listing
	listing, err := act.CreateListing(sellerID, itemID, 2, 20)
	if err != nil {
		t.Fatalf("Erreur CreateListing: %v", err)
	}
	
	// Etape 2 : Achat (Transaction)
	err = act.BuyListing(buyerID, listing.CustomID)
	if err != nil {
		t.Fatalf("Erreur BuyListing: %v", err)
	}
	
	// Etape 3 : Verifications POST achat
	var buyer models.Player
	playerCol.FindOne(context.TODO(), map[string]string{"customID": buyerID}).Decode(&buyer)
	if buyer.Gold != 60 { // 100 - (2*20) = 60
		t.Errorf("Le buyer devrait avoir 60 gold, il en a %d", buyer.Gold)
	}
	
	var seller models.Player
	playerCol.FindOne(context.TODO(), map[string]string{"customID": sellerID}).Decode(&seller)
	if seller.Gold != 40 { // 0 + (2*20) = 40
		t.Errorf("Le seller devrait avoir 40 gold, il en a %d", seller.Gold)
	}
	
	var buyerInv models.InventoryEntry
	invCol.FindOne(context.TODO(), map[string]string{"playerID": buyerID, "itemID": itemID}).Decode(&buyerInv)
	if buyerInv.Qty != 2 {
		t.Errorf("Le buyer devrait avoir 2 items, il en a %d", buyerInv.Qty)
	}
	
	// Etape 4 : Double Achat impossible
	err = act.BuyListing(buyerID, listing.CustomID)
	if err == nil {
		t.Fatalf("Le double achat devrait renvoyer une erreur car le statut n'est plus active")
	}
}
