package internal

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const collection = "todos"

type MongoHandler struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHandler(address string, database string) *MongoHandler {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cl, _ := mongo.Connect(ctx, options.Client().ApplyURI(address))
	coll := cl.Database(database).Collection(collection)
	return &MongoHandler{
		coll: coll,
	}
}

func (mh *MongoHandler) GetOne(c *Todo, filter interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := mh.coll.FindOne(ctx, filter).Decode(c)
	return err
}

func (mh *MongoHandler) Get(filter interface{}) []*Todo {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur, err := mh.coll.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	var result []*Todo
	for cur.Next(ctx) {
		contact := &Todo{}
		er := cur.Decode(contact)
		if er != nil {
			log.Fatal(er)
		}
		result = append(result, contact)
	}
	return result
}

func (mh *MongoHandler) AddOne(c *Todo) (*mongo.InsertOneResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := mh.coll.InsertOne(ctx, c)
	return result, err
}

func (mh *MongoHandler) Update(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := mh.coll.UpdateMany(ctx, filter, update)
	return result, err
}

func (mh *MongoHandler) RemoveOne(filter interface{}) (*mongo.DeleteResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := mh.coll.DeleteOne(ctx, filter)
	return result, err
}
