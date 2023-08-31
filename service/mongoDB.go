package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Acc is
type Acc struct {
	ID     string `bson:"_id,omitempty" json:"uid"`
	Acc    string `bson:"acc" json:"acc"`
	Pswd   string `bson:"pswd" json:"pswd"`
	Email  string `bson:"email" json:"email"`
	Name   string `bson:"name" json:"name"`
	Gender string `bson:"gender" json:"gender"`
}

// ConnectionInfo is
type ConnectionInfo struct {
	DBName         string
	CollectionName string
}

// MongoDBcontext is connect setting
func (c ConnectionInfo) MongoDBcontext() *mongo.Collection {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pswd := os.Getenv("MONGODBPSWD")

	conn := fmt.Sprint("mongodb+srv://j_dev:", pswd, "@jdev.y4x5s.gcp.mongodb.net/?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conn))
	if err != nil {
		log.Fatal(err)
	}

	return client.Database(c.DBName).Collection(c.CollectionName)
}
