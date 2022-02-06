package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Price struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Name  string             `json:"name" bson:"name"`
	Price float32            `json:"price" bson:"price"`
}
