package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    Id       			   primitive.ObjectID 	`json:"id,omitempty"`
	DNI             	   string         		`json:"dni,omitempty"`
	FirstName       	   string         		`json:"firstName,omitempty"`
	LastName        	   string         		`json:"lastName,omitempty"`
	CivilStatus     	   string         		`json:"civilStatus,omitempty"`
	BirthDateString        string              	`json:"birthDateString,omitempty" bson:"-"`
	BirthDate              time.Time           	`json:"birthDate,omitempty" bson:"birthDate,omitempty"`
	City            	   string            	`json:"city,omitempty"`
	CityID            	   string            	`json:"cityId,omitempty"`
	Email           	   string         		`json:"email,omitempty"`
	Phone           	   string         		`json:"phone,omitempty"`
}