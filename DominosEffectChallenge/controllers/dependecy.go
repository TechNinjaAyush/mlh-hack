package controllers


import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	Mongo *mongo.Database
}

func NewHandler(mongo *mongo.Database) *Handler {
	return &Handler{
		Mongo: mongo,
	}
}
