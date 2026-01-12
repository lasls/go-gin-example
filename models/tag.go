package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Tag struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	State      int                `bson:"state" json:"state"`
	CreatedBy  string             `bson:"created_by" json:"created_by"`
	ModifiedBy string             `bson:"modified_by" json:"modified_by"`
	CreatedOn  time.Time          `bson:"created_on" json:"created_on"`
	ModifiedOn time.Time          `bson:"modified_on" json:"modified_on"`
}

const (
	CollectionTag = "tags"
)

func GetTags(pageNum int, pageSize int, maps bson.M) ([]Tag, error) {
	ctx := context.Background()
	findOptions := options.Find()
	findOptions.SetSkip(int64(pageNum))
	findOptions.SetLimit(int64(pageSize))

	cursor, err := db.Collection(CollectionTag).Find(ctx, maps, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tags []Tag
	if err = cursor.All(ctx, &tags); err != nil {
		return nil, err
	}

	return tags, nil
}

func GetTagTotal(maps bson.M) (int64, error) {
	ctx := context.Background()
	count, err := db.Collection(CollectionTag).CountDocuments(ctx, maps)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func AddTag(name string, state int, createdBy string) error {
	ctx := context.Background()
	tag := Tag{
		Name:       name,
		State:      state,
		CreatedBy:  createdBy,
		CreatedOn:  time.Now(),
		ModifiedOn: time.Now(),
	}
	_, err := db.Collection(CollectionTag).InsertOne(ctx, tag)
	return err
}

func EditTag(id string, data bson.M) error {
	ctx := context.Background()
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	data["modified_on"] = time.Now()
	_, err = db.Collection(CollectionTag).UpdateByID(ctx, objectId, bson.M{"$set": data})
	return err
}

func DeleteTag(id string) error {
	ctx := context.Background()
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = db.Collection(CollectionTag).DeleteOne(ctx, bson.M{"_id": objectId})
	return err
}

func ExistTagByName(name string) (bool, error) {
	ctx := context.Background()
	count, err := db.Collection(CollectionTag).CountDocuments(ctx, bson.M{"name": name})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func ExistTagByID(id string) (bool, error) {
	ctx := context.Background()
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}
	count, err := db.Collection(CollectionTag).CountDocuments(ctx, bson.M{"_id": objectId})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
