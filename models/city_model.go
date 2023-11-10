package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type City struct {
    Key       			   primitive.ObjectID 	`json:"key,omitempty"`
	Value             	   string         		`json:"value,omitempty"`
}