package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://suryansh:suru123@cluster0.eik0l.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

var UserCollection *mongo.Collection
var ChatCollection *mongo.Collection

func init() {

	// client option
	clientOption := options.Client().ApplyURI(connectionString)

	//connect to mongoDB
	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}

	UserCollection = client.Database("vartalap").Collection("userCollection")
	ChatCollection = client.Database("vartalap").Collection("chatCollection")

	fmt.Println("Vartalap database is ready to use")
}
