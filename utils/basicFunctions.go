package utils

import (
	"context"
	"vartalap/database"
	"vartalap/model"

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
		temp, _ := GetAnotherUser(user.Chatsid[i], userID)
		chatedUsers = append(chatedUsers, temp)
	}
	return chatedUsers

}

func GetAnotherUser(user1 primitive.ObjectID, user2 string) (primitive.ObjectID, error) {
	var chat model.Chat
	err := database.ChatCollection.FindOne(context.TODO(), bson.M{
		"_id": user1,
	}).Decode(&chat)
	if err != nil {
		return user1, err
	}
	if chat.User1.String() == user2 {
		return chat.User2, nil
	} else {
		return chat.User1, nil
	}
}
