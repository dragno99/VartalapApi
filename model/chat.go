package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	User1    primitive.ObjectID `bson:"user1,omitempty"`
	User2    primitive.ObjectID `bson:"user2,omitempty"`
	Messages []Message          `bson:"messages,omitempty"`
}

type Message struct {
	Text      string             `json:"text"`
	Image     string             `json:"image"`
	SenderId  primitive.ObjectID `json:"senderId"`
	CreatedAt int64              `json:"createdAt"`
}
