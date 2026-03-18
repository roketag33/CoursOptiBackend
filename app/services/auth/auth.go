package auth

import (
	"context"
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/server"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct{}

func New() *Auth {
	return &Auth{}
}

func (a *Auth) Register(displayName, password string) (*models.Player, error) {
	srv := server.GetServer()
	col := srv.Database.Collection("player")

	var existing models.Player
	err := col.FindOne(context.TODO(), bson.M{"displayName": displayName}).Decode(&existing)
	if err == nil {
		return nil, errors.New("ce nom d'utilisateur est deja pris")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	player := models.Player{
		CustomID:    functions.NewUUID(),
		DisplayName: displayName,
		Password:    string(hashedPassword),
		Gold:        0,
		Stats:       map[string]interface{}{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = col.InsertOne(context.TODO(), player)
	if err != nil {
		return nil, err
	}

	return &player, nil
}

func (a *Auth) Login(displayName, password string) (string, error) {
	srv := server.GetServer()
	col := srv.Database.Collection("player")

	var player models.Player
	err := col.FindOne(context.TODO(), bson.M{"displayName": displayName}).Decode(&player)
	if err != nil {
		return "", errors.New("identifiants invalides")
	}

	err = bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(password))
	if err != nil {
		return "", errors.New("identifiants invalides")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"customID": player.CustomID,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
