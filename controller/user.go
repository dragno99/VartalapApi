package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dragno99/VartalapApi/database"
	"github.com/dragno99/VartalapApi/model"
	"github.com/dragno99/VartalapApi/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// function for getting user chats
func GetUserChats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId, _ := primitive.ObjectIDFromHex(r.Header.Get("userId"))
	var user model.User
	err := database.UserCollection.FindOne(context.TODO(), bson.M{
		"_id": userId,
	}).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	type tempData struct {
		Fullname string             `json:"fullname"`
		Username string             `json:"username"`
		ID       primitive.ObjectID `json:"_id"`
		Imageurl string             `json:"imageurl"`
		Chatid   primitive.ObjectID `json:"chatid"`
		Pubkey   string             `json:"pubkey"`
	}

	var userChats = make([]tempData, 0)
	for i := 0; i < len(user.Chatsid); i++ {
		anotherUserId, err := utils.GetAnotherUser(user.Chatsid[i], r.Header.Get("userId"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var chattedUser model.User
		err = database.UserCollection.FindOne(context.TODO(), bson.M{
			"_id": anotherUserId,
		}).Decode(&chattedUser)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userChats = append(userChats, tempData{
			Username: chattedUser.Username,
			ID:       chattedUser.ID,
			Fullname: chattedUser.Fullname,
			Chatid:   user.Chatsid[i],
			Imageurl: chattedUser.Imageurl,
			Pubkey:   chattedUser.Pubkey,
		})
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(bson.M{
		"users": userChats,
	})

}

// function for getting messages related to a particular chat id
func GetChatMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	if params["chatId"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bson.M{
			"message": "Please provide a chatId",
		})
		return
	}
	chatId := params["chatId"]
	var chat model.Chat
	id, err := primitive.ObjectIDFromHex(chatId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bson.M{
			"message": "chatId is not valid",
		})
		return
	}

	if err := database.ChatCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&chat); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bson.M{
			"message": "chatId is not valid",
		})
		return
	}
	w.WriteHeader(http.StatusAccepted)

	if chat.Messages == nil {
		chat.Messages = make([]model.Message, 0)
	}

	json.NewEncoder(w).Encode(bson.M{
		"messages": chat.Messages,
	})
}

// function for seeing the list of others users in the database
func GetAppUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	chatedUsers := utils.GetChatedUsers(r.Header.Get("userId"))
	userId, _ := primitive.ObjectIDFromHex(r.Header.Get("userId"))
	chatedUsers = append(chatedUsers, userId)
	filter := bson.M{
		"_id": bson.M{
			"$nin": chatedUsers,
		},
	}
	cursor, err := database.UserCollection.Find(context.TODO(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	type userData struct {
		Fullname string             `json:"fullname"`
		Username string             `json:"username"`
		Id       primitive.ObjectID `json:"_id"`
		Imageurl string             `json:"imageurl"`
	}

	var resUser = make([]userData, 0)

	for cursor.Next(context.TODO()) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resUser = append(resUser, userData{
			Fullname: user.Fullname,
			Username: user.Username,
			Id:       user.ID,
			Imageurl: user.Imageurl,
		})
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(bson.M{
		"users": resUser,
	})
}

// function to start chat with a particular user
func StartChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// params := mux.Vars(r)
	userId, err := primitive.ObjectIDFromHex(r.Header.Get("userId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	type tempData struct {
		UserId primitive.ObjectID `json:"userId"`
	}
	var user2 tempData
	err = json.Unmarshal(body, &user2)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user2id := user2.UserId
	var chat model.Chat
	chat.User1 = userId
	chat.User2 = user2id
	insertedChat, err := database.ChatCollection.InsertOne(context.TODO(), chat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pushChatIntoUserChatArray := func(userId primitive.ObjectID, chatId primitive.ObjectID) bool {
		filter := bson.M{
			"_id": userId,
		}
		update := bson.M{
			"$push": bson.M{
				"chatsid": chatId,
			},
		}
		_, err = database.UserCollection.UpdateOne(context.TODO(), filter, update)
		return err == nil
	}
	bool1 := pushChatIntoUserChatArray(userId, insertedChat.InsertedID.(primitive.ObjectID))
	bool2 := pushChatIntoUserChatArray(user2id, insertedChat.InsertedID.(primitive.ObjectID))
	if !bool1 || !bool2 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(bson.M{
		"message": "Started chat successfully",
	})
}

// function to add chat message in a particular chatId
func AddMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// params := mux.Vars(r)
	userId, err := primitive.ObjectIDFromHex(r.Header.Get("userId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	type MsgBody struct {
		Text      string             `json:"text,omitempty"`
		Image     string             `json:"image,omitempty"`
		CreatedAt int64              `json:"createdAt,omitempty"`
		ChatId    primitive.ObjectID `json:"chatId,omitempty"`
	}

	var msg MsgBody
	err = json.Unmarshal(body, &msg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	currMessage := model.Message{
		Text:      msg.Text,
		Image:     msg.Image,
		SenderId:  userId,
		CreatedAt: msg.CreatedAt,
	}

	update := bson.M{
		"$push": bson.M{
			"messages": currMessage,
		},
	}
	result, err := database.ChatCollection.UpdateByID(context.TODO(), msg.ChatId, update)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if result.MatchedCount > 0 {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(bson.M{
			"message": "message added successfully.",
		})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bson.M{
			"message": "message sending failed",
		})
	}

}

// function to update user's fullname
func UpdateFullName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// params := mux.Vars(r)
	userId, err := primitive.ObjectIDFromHex(r.Header.Get("userId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bson.M{
			"message": "please provide a valid userId(user object id)",
		})
		return
	}
	type tempData struct {
		UpdatedFullname string `json:"updatedFullname,omitempty"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	var temp tempData
	json.Unmarshal(body, &temp)
	filter := bson.M{
		"_id": userId,
	}
	update := bson.M{
		"$set": bson.M{
			"fullname": temp.UpdatedFullname,
		},
	}
	_, err = database.UserCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(bson.M{
		"message": "user's fullname updated",
	})
}
