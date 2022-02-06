package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Username string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
}
