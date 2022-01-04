package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"vartalap/database"
	"vartalap/model"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var userCollection = database.UserCollection

func UserSignUp(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(bson.M{
			"message": "hash nai bna",
		})
		return
	}
	user.Password = string(hashedPassword)
	_, err = userCollection.InsertOne(context.TODO(), user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(bson.M{
			"message": "hash to hua but store nahi hua db mein",
		})
		return
	}
	json.NewEncoder(w).Encode(bson.M{
		"message": "User Signup Successfully",
	})
}

func UserLogIn(w http.ResponseWriter, r *http.Request) {
	var user, foundUser model.User
	json.NewDecoder(r.Body).Decode(&user)
	filter := bson.M{"username": user.Username}
	cursor, err := userCollection.Find(context.TODO(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for cursor.Next(context.TODO()) {
		err = cursor.Decode(&foundUser)
		if err != nil {
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(bson.M{
				"message": "No user Found",
			})
			return
		} else if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(bson.M{
				"message": "No user Found",
			})
			return
		} else {
			cursor.Close(context.TODO())
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(bson.M{
				"message": "User login Successfully",
			})
			return
		}
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(bson.M{
		"message": "No user Found",
	})
}

func SayHello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(bson.M{
		"message": "hello from suru chat api",
	})
}
