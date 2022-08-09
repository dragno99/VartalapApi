package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dragno99/VartalapApi/database"
	"github.com/dragno99/VartalapApi/middleware"
	"github.com/dragno99/VartalapApi/model"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var userCollection = database.UserCollection

func UserSignUp(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bson.M{
			"message": "Please provide a data",
		})
		return
	}

	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user, err = formatAndValidateForSignUp(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if alreadyExists(user.Username) {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(bson.M{
			"message": "Username already exists",
		})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)
	_, err = userCollection.InsertOne(context.TODO(), user)

	if err != nil {

		fmt.Println("yaha", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(bson.M{
		"message": "User Signup Successfully",
	})
}

func UserLogIn(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bson.M{
			"message": "Please provide a data",
		})
		return
	}

	var inputData, foundUser model.User
	err := json.NewDecoder(r.Body).Decode(&inputData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	inputData, err = formatAndValidateForLogIn(inputData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bson.M{
			"message": "Bad input data",
		})
		return
	}
	filter := bson.M{"username": inputData.Username}
	if err := userCollection.FindOne(context.TODO(), filter).Decode(&foundUser); err != nil {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(bson.M{
			"message": "No user Found",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(inputData.Password)); err != nil {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(bson.M{
			"message": "No user Found",
		})
		return
	}

	// update pubkey
	_, err = userCollection.UpdateByID(context.TODO(), foundUser.ID, bson.M{
		"$set": bson.M{
			"pubkey": inputData.Pubkey,
		},
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(bson.M{
			"message": "User login failed",
		})
	}

	//user login successfully
	// create JWT
	tokenString, err := middleware.GenerateJWT(foundUser.ID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(bson.M{
			"message": "User login failed",
		})
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(bson.M{
		"message":  "User login Successfully",
		"token":    tokenString,
		"userId":   foundUser.ID,
		"fullname": foundUser.Fullname,
		"imageurl": foundUser.Imageurl,
	})
}

func SayHello(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(bson.M{
		"message": "hello from Vartalap chat api",
	})
}

func formatAndValidateForSignUp(user model.User) (model.User, error) {
	if len(user.Username) == 0 {
		return user, errors.New("please provide a username")
	} else {
		user.Username = strings.TrimSpace(user.Username)
	}

	if len(user.Password) == 0 {
		return user, errors.New("please provide a password")
	} else {
		user.Password = strings.TrimSpace(user.Password)
	}

	if len(user.Email) == 0 {
		return user, errors.New("please provide a email")
	} else {
		user.Email = strings.TrimSpace(user.Email)
	}

	if len(user.Pubkey) == 0 {
		return user, errors.New("please provide a pubkey")
	} else {
		user.Pubkey = strings.TrimSpace(user.Pubkey)
	}
	if len(user.Imagename) > 0 {
		user.Imagename = strings.TrimSpace(user.Imagename)
	}
	if len(user.Imageurl) > 0 {
		user.Imageurl = strings.TrimSpace(user.Imageurl)
	}
	return user, nil
}

func formatAndValidateForLogIn(user model.User) (model.User, error) {
	if len(user.Username) == 0 {
		return user, errors.New("please provide a username")
	} else {
		user.Username = strings.TrimSpace(user.Username)
	}
	if len(user.Password) == 0 {
		return user, errors.New("please provide a password")
	} else {
		user.Password = strings.TrimSpace(user.Password)
	}

	if len(user.Pubkey) == 0 {
		return user, errors.New("please provide a pubkey")
	} else {
		user.Pubkey = strings.TrimSpace(user.Pubkey)
	}
	return user, nil
}

func alreadyExists(username string) bool {
	var tempUser model.User
	_ = userCollection.FindOne(context.TODO(), bson.M{
		"username": username,
	}).Decode(&tempUser)
	return len(tempUser.Username) > 0
}
