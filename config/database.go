package config

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	mongoString := os.Getenv("MONGOSTRING") // ambil dari .env di sini
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoString))
	if err != nil {
		panic(err)
	}
	DB = client.Database("dbFilm")
	fmt.Println("MongoDB Connected")
}
