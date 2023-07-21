package service

import (
	"context"
	"fmt"
	"log"
	"os"
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
	conn := fmt.Sprint("mongodb+srv://j_dev:", os.Getenv("MONGODBPSWD"), "@jdev.y4x5s.gcp.mongodb.net/?retryWrites=true&w=majority")
	fmt.Println("get env")
	fmt.Println(os.Getenv("MONGODBPSWD"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conn))
	if err != nil {
		log.Fatal(err)
	}

	return client.Database(c.DBName).Collection(c.CollectionName)
}
