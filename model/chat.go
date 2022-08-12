package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	User1    primitive.ObjectID `bson:"user1,omitempty"`
	User2    primitive.ObjectID `bson:"user2,omitempty"`
	Messages []Message          `bson:"messages"`
}

type Message struct {
	Text      string             `bson:"text"`
	Image     string             `bson:"image"`
	SenderId  primitive.ObjectID `bson:"senderId"`
	CreatedAt int64              `bson:"createdAt"`
}
