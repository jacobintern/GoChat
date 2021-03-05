package service

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Acc is
type Acc struct {
	ID     string `bson:"_id,omitempty"`
	Acc    string `bson:"acc"`
	Pswd   string `bson:"pswd"`
	Email  string `bson:"email"`
	Name   string `bson:"name"`
	Gender string `bson:"gender"`
}

// ConnectionInfo is
type ConnectionInfo struct {
	DBName         string
	CollectionName string
}

// MongoDBcontext is connect setting
func (c ConnectionInfo) MongoDBcontext() *mongo.Collection {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://j_dev:zHYJQ2jc7UAqHThV@jdev.y4x5s.gcp.mongodb.net/"+c.DBName+"?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}
	return client.Database(c.DBName).Collection(c.CollectionName)
}
