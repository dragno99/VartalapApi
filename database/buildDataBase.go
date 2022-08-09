package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection *mongo.Collection
var ChatCollection *mongo.Collection

func init() {

	// client option
	clientOption := options.Client().ApplyURI(getConnectionString())

	//connect to mongoDB
	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}

	UserCollection = client.Database("vartalap").Collection("userCollection")
	ChatCollection = client.Database("vartalap").Collection("chatCollection")

	fmt.Println("Vartalap database is ready to use...")
}

func getConnectionString() string {
	godotenv.Load()
	username := os.Getenv("USERNAME_mongo")
	password := os.Getenv("PASSWORD_mongo")
	if username == "" {
		username = "suryansh"
	}
	if password == "" {
		password = "suru123"
	}
	connectionString := "mongodb+srv://" + username + ":" + password + "@cluster0.lw6zntr.mongodb.net/?retryWrites=true&w=majority"
	return connectionString
}
