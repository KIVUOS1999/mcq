package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type QuestionStruct struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Question string             `json:"question"`
	Options  map[int]string     `json:"options"`
	Answer   string             `json:"answer"`
}
