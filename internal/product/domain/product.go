package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID       string             `json:"id" bson:"-"`
	ObjectID primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Name     string             `json:"name" bson:"name"`
	Category string             `json:"category" bson:"category"`
	Price    float64            `json:"price" bson:"price"`
	Stock    int                `json:"stock" bson:"stock"`
	Brand    string             `json:"brand" bson:"brand"`
}
