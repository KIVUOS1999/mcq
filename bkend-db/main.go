package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/bkend-db/handler"
	"github.com/bkend-db/store"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load("configs/.env")
	if err != nil {
		panic("Error loading .env file")
	}

	cli, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_ENDPOINT")))
	if err != nil {
		log.Fatalf(err.Error())
	}

	defer cli.Disconnect(context.Background())

	s := store.New(cli)
	handler := handler.New(s)

	router := mux.NewRouter()
	router.HandleFunc("/get_question/{count}", handler.GetCountedRecords).Methods(http.MethodGet)
	router.HandleFunc("/get_question/id/{id}", handler.GetQuestionByID).Methods(http.MethodGet)

	http.ListenAndServe(os.Getenv("HTTP_PORT"), router)
}
