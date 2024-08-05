package store

import (
	"context"
	"log"

	"github.com/bkend-db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	collection *mongo.Collection
}

func New(client *mongo.Client) *Store {
	err := client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("connected to mongo !!!!")

	collection := client.Database("mcq").Collection("questions")

	return &Store{collection: collection}
}

func (s *Store) GetRecords(count int) ([]models.QuestionStruct, error) {
	ctx := context.Background()

	pipeline := bson.A{
		bson.D{{"$sample", bson.D{{"size", count}}}},
	}

	cur, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Print("[ERROR] -", err.Error())

		return nil, err
	}

	defer cur.Close(ctx)

	var questionSet []models.QuestionStruct

	for cur.Next(ctx) {
		var question models.QuestionStruct

		err := cur.Decode(&question)
		if err != nil {
			log.Print("[ERROR] -", err)

			return nil, err
		}

		questionSet = append(questionSet, question)
	}

	if err := cur.Err(); err != nil {
		log.Print("[ERROR] -", err)

		return nil, err
	}

	return questionSet, nil
}

func (s *Store) GetQuestionByID(id string) (*models.QuestionStruct, error) {
	ctx := context.Background()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("[ERROR] :", err.Error())

		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var question models.QuestionStruct

	err = s.collection.FindOne(ctx, filter).Decode(&question)
	if err != nil {
		log.Println(err)

		return nil, err
	}

	return &question, nil
}
