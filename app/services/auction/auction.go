package auction

import (
	"context"
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/server"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Auction struct{}

func New() *Auction {
	return &Auction{}
}

func (a *Auction) CreateListing(sellerID string, itemID string, qty int, pricePerUnit int64) (*models.Listing, error) {
	srv := server.GetServer()

	listing := models.Listing{
		CustomID:     functions.NewUUID(),
		SellerID:     sellerID,
		ItemID:       itemID,
		Qty:          qty,
		PricePerUnit: pricePerUnit,
		Status:       "active",
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(24 * 7 * time.Hour),
	}

	session, err := srv.Database.Client().StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(context.TODO(), func(sessCtx context.Context) (interface{}, error) {
		invCol := srv.Database.Collection("inventory")
		listingCol := srv.Database.Collection("listing")

		var inv models.InventoryEntry
		err := invCol.FindOne(sessCtx, bson.M{"playerID": sellerID, "itemID": itemID}).Decode(&inv)
		if err != nil || inv.Qty < int64(qty) {
			return nil, errors.New("pas assez d'items en inventaire")
		}

		_, err = invCol.UpdateOne(sessCtx, bson.M{"playerID": sellerID, "itemID": itemID}, bson.M{
			"$inc": bson.M{"qty": -qty},
			"$set": bson.M{"updatedAt": time.Now()},
		})
		if err != nil {
			return nil, err
		}

		_, err = listingCol.InsertOne(sessCtx, listing)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return nil, err
	}

	return &listing, nil
}

func (a *Auction) GetActiveListings() ([]models.Listing, error) {
	srv := server.GetServer()
	col := srv.Database.Collection("listing")

	var listings []models.Listing
	cursor, err := col.Find(context.TODO(), bson.M{"status": "active"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var l models.Listing
		if err := cursor.Decode(&l); err != nil {
			return nil, err
		}
		listings = append(listings, l)
	}

	return listings, nil
}

func (a *Auction) BuyListing(buyerID string, listingID string) error {
	srv := server.GetServer()

	session, err := srv.Database.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(context.TODO(), func(sessCtx context.Context) (interface{}, error) {
		listingCol := srv.Database.Collection("listing")
		playerCol := srv.Database.Collection("player")
		invCol := srv.Database.Collection("inventory")

		var listing models.Listing
		err := listingCol.FindOne(sessCtx, bson.M{"customID": listingID}).Decode(&listing)
		if err != nil {
			return nil, errors.New("listing introuvable")
		}

		if listing.Status != "active" {
			return nil, errors.New("listing n'est plus actif")
		}

		if listing.SellerID == buyerID {
			return nil, errors.New("tu ne peux pas acheter ton propre item")
		}

		totalPrice := listing.PricePerUnit * int64(listing.Qty)

		var buyer models.Player
		err = playerCol.FindOne(sessCtx, bson.M{"customID": buyerID}).Decode(&buyer)
		if err != nil || buyer.Gold < totalPrice {
			return nil, errors.New("fonds insuffisants")
		}

		// Debit buyer
		_, err = playerCol.UpdateOne(sessCtx, bson.M{"customID": buyerID}, bson.M{"$inc": bson.M{"gold": -totalPrice}})
		if err != nil {
			return nil, err
		}

		// Credit seller
		_, err = playerCol.UpdateOne(sessCtx, bson.M{"customID": listing.SellerID}, bson.M{"$inc": bson.M{"gold": totalPrice}})
		if err != nil {
			return nil, err
		}

		// Transfer item
		var existing models.InventoryEntry
		err = invCol.FindOne(sessCtx, bson.M{"playerID": buyerID, "itemID": listing.ItemID}).Decode(&existing)
		if err != nil {
			_, err = invCol.InsertOne(sessCtx, models.InventoryEntry{
				PlayerID:  buyerID,
				ItemID:    listing.ItemID,
				Qty:       int64(listing.Qty),
				UpdatedAt: time.Now(),
			})
		} else {
			_, err = invCol.UpdateOne(sessCtx, bson.M{"playerID": buyerID, "itemID": listing.ItemID}, bson.M{
				"$inc": bson.M{"qty": listing.Qty},
				"$set": bson.M{"updatedAt": time.Now()},
			})
		}
		if err != nil {
			return nil, err
		}

		// Mark sold
		_, err = listingCol.UpdateOne(sessCtx, bson.M{"customID": listingID}, bson.M{"$set": bson.M{"status": "sold"}})
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

func (a *Auction) CancelListing(sellerID string, listingID string) error {
	srv := server.GetServer()

	session, err := srv.Database.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(context.TODO(), func(sessCtx context.Context) (interface{}, error) {
		listingCol := srv.Database.Collection("listing")
		invCol := srv.Database.Collection("inventory")

		var listing models.Listing
		err := listingCol.FindOne(sessCtx, bson.M{"customID": listingID}).Decode(&listing)
		if err != nil {
			return nil, errors.New("listing introuvable")
		}

		if listing.Status != "active" {
			return nil, errors.New("listing n'est plus actif")
		}

		if listing.SellerID != sellerID {
			return nil, errors.New("ce n'est pas ton listing")
		}

		// Rend l'item
		_, err = invCol.UpdateOne(sessCtx, bson.M{"playerID": sellerID, "itemID": listing.ItemID}, bson.M{
			"$inc": bson.M{"qty": listing.Qty},
			"$set": bson.M{"updatedAt": time.Now()},
		})
		if err != nil {
			return nil, err
		}

		// Mark cancelled
		_, err = listingCol.UpdateOne(sessCtx, bson.M{"customID": listingID}, bson.M{"$set": bson.M{"status": "cancelled"}})
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}
