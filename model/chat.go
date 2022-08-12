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
	Text      string             `bson:"text,omitempty"`
	Image     string             `bson:"image,omitempty"`
	SenderId  primitive.ObjectID `bson:"senderId,omitempty"`
	CreatedAt int64              `bson:"createdAt,omitempty"`
}
