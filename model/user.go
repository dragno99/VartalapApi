package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	Username  string               `bson:"username,omitempty"`
	Fullname  string               `bson:"fullname,omitempty"`
	Imageurl  string               `bson:"imageurl,omitempty"`
	Imagename string               `bson:"imagename,omitempty"`
	Pubkey    string               `bson:"pubkey,omitempty"`
	Password  string               `bson:"password,omitempty"`
	Email     string               `bson:"email,omitempty"`
	Chatsid   []primitive.ObjectID `bson:"chatsid,omitempty"`
}
