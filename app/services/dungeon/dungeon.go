package dungeon

import (
	"context"
	"dungeons/app/functions"
	"dungeons/app/models"
	"dungeons/app/mongodb"
	"dungeons/app/server"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Dungeon struct{}

func New() *Dungeon {
	return &Dungeon{}
}

func (d *Dungeon) Create(in *models.Dungeon) (*models.Dungeon, error) {
	srv := server.GetServer()
	col := srv.Database.Collection("dungeon")

	in.CustomID = functions.NewUUID()
	in.Status = "draft"
	in.CreatedAt = time.Now()
	in.UpdatedAt = time.Now()
	if in.Bosses == nil {
		in.Bosses = []string{}
	}

	_, err := col.InsertOne(context.TODO(), in)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func (d *Dungeon) GetByID(id string) (*models.Dungeon, error) {
	srv := server.GetServer()
	col := srv.Database.Collection("dungeon")

	var dungeon models.Dungeon
	filter := bson.M{"customID": id}
	err := col.FindOne(context.TODO(), filter).Decode(&dungeon)
	if err != nil {
		return nil, err
	}

	return &dungeon, nil
}

func (d *Dungeon) GetAll() ([]models.Dungeon, error) {
	srv := server.GetServer()
	col := srv.Database.Collection("dungeon")

	var dungeons []models.Dungeon
	cursor, err := col.Find(context.TODO(), bson.M{"status": "published"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var dg models.Dungeon
		if err := cursor.Decode(&dg); err != nil {
			return nil, err
		}
		dungeons = append(dungeons, dg)
	}

	return dungeons, nil
}

func (d *Dungeon) Update(id string, in *models.Dungeon) error {
	srv := server.GetServer()
	col := srv.Database.Collection("dungeon")

	existing, err := d.GetByID(id)
	if err != nil {
		return err
	}

	if existing.Status == "published" {
		return errors.New("impossible de modifier un dungeon publié")
	}

	in.UpdatedAt = time.Now()

	doc, err := mongodb.ToDoc(in)
	if err != nil {
		return err
	}
	delete(doc, "customID")
	delete(doc, "createdAt")

	filter := bson.M{"customID": id}
	update := bson.M{"$set": doc}
	_, err = col.UpdateOne(context.TODO(), filter, update)
	return err
}

func (d *Dungeon) Publish(id string) error {
	srv := server.GetServer()
	col := srv.Database.Collection("dungeon")

	dungeon, err := d.GetByID(id)
	if err != nil {
		return err
	}

	if len(dungeon.Bosses) == 0 {
		return errors.New("impossible de publier un dungeon sans boss")
	}

	filter := bson.M{"customID": id}
	update := bson.M{"$set": bson.M{"status": "published", "updatedAt": time.Now()}}
	_, err = col.UpdateOne(context.TODO(), filter, update)
	return err
}

func (d *Dungeon) AddStep(dungeonID string, in *models.BossStep) (*models.BossStep, error) {
	srv := server.GetServer()
	bossCol := srv.Database.Collection("bossstep")
	dungeonCol := srv.Database.Collection("dungeon")

	dungeon, err := d.GetByID(dungeonID)
	if err != nil {
		return nil, errors.New("dungeon pas trouvé")
	}

	in.CustomID = functions.NewUUID()
	in.DungeonID = dungeonID
	in.Order = len(dungeon.Bosses) + 1

	_, err = bossCol.InsertOne(context.TODO(), in)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"customID": dungeonID}
	update := bson.M{
		"$push": bson.M{"bosses": in.CustomID},
		"$set":  bson.M{"updatedAt": time.Now()},
	}
	_, err = dungeonCol.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func (d *Dungeon) UpdateStep(dungeonID string, stepID string, in *models.BossStep) error {
	srv := server.GetServer()
	col := srv.Database.Collection("bossstep")

	_, err := d.GetByID(dungeonID)
	if err != nil {
		return errors.New("dungeon pas trouvé")
	}

	doc, err := mongodb.ToDoc(in)
	if err != nil {
		return err
	}
	delete(doc, "customID")
	delete(doc, "dungeonID")

	filter := bson.M{"customID": stepID, "dungeonID": dungeonID}
	update := bson.M{"$set": doc}
	result, err := col.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("step pas trouvé dans ce dungeon")
	}

	return nil
}

func (d *Dungeon) ReorderSteps(dungeonID string, stepIDs []string) error {
	srv := server.GetServer()
	bossCol := srv.Database.Collection("bossstep")
	dungeonCol := srv.Database.Collection("dungeon")

	_, err := d.GetByID(dungeonID)
	if err != nil {
		return errors.New("dungeon pas trouvé")
	}

	for i, stepID := range stepIDs {
		filter := bson.M{"customID": stepID, "dungeonID": dungeonID}
		update := bson.M{"$set": bson.M{"order": i + 1}}
		_, err := bossCol.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return err
		}
	}

	filter := bson.M{"customID": dungeonID}
	update := bson.M{
		"$set": bson.M{"bosses": stepIDs, "updatedAt": time.Now()},
	}
	_, err = dungeonCol.UpdateOne(context.TODO(), filter, update)
	return err
}

func (d *Dungeon) GetSteps(dungeonID string) ([]models.BossStep, error) {
	srv := server.GetServer()
	col := srv.Database.Collection("bossstep")

	var steps []models.BossStep
	filter := bson.M{"dungeonID": dungeonID}

	cursor, err := col.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var step models.BossStep
		if err := cursor.Decode(&step); err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}

	return steps, nil
}
