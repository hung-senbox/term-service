package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Term struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Title            string             `bson:"title"`
	Color            string             `bson:"color"`
	PublishedMobile  bool               `bson:"published_mobile"`
	PublishedDesktop bool               `bson:"published_desktop"`
	StartDate        time.Time          `bson:"start_date"`
	EndDate          time.Time          `bson:"end_date"`
	CreatedAt        time.Time          `bson:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at"`
}
