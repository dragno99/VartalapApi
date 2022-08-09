package utils

import (
	"context"
	"log"

	"github.com/dragno99/vartalapAPI/database"
	"github.com/dragno99/vartalapAPI/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetChatedUsers(userID string) []primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(userID)
	var user model.User
	database.UserCollection.FindOne(context.TODO(), bson.M{
		"_id": id,
	}).Decode(&user)
	var chatedUsers []primitive.ObjectID
	for i := 0; i < len(user.Chatsid); i++ {
		temp, err := GetAnotherUser(user.Chatsid[i], userID)
		if err != nil {
			log.Fatal(err)
			return chatedUsers
		}
		chatedUsers = append(chatedUsers, temp)
	}
	return chatedUsers

}

func GetAnotherUser(chatId primitive.ObjectID, user string) (primitive.ObjectID, error) {
	var chat model.Chat
	err := database.ChatCollection.FindOne(context.TODO(), bson.M{
		"_id": chatId,
	}).Decode(&chat)
	if err != nil {
		return chatId, err
	}
	if chat.User1.Hex() == user {
		return chat.User2, nil
	} else {
		return chat.User1, nil
	}
}
