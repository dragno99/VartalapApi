package wSocket

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	Text      string             `json:"text"`
	Image     string             `json:"image"`
	SenderId  primitive.ObjectID `json:"senderId"`
	CreatedAt int64              `json:"createdAt"`
}
