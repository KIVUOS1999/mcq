package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mcq_backend/constants"
	"github.com/mcq_backend/handler"
	"github.com/mcq_backend/storage"

	dataClient "github.com/bkend-db/pkg/client"
	redisClient "github.com/bkend-redis/pkg/client"
)

func main() {
	err := godotenv.Load("configs/.env")
	if err != nil {
		panic("Error loading .env file")
	}

	router := mux.NewRouter()

	storage := storage.NewStorage()
	redisClient := redisClient.New(os.Getenv("REDIS_SVC"))
	dataClient := dataClient.New(os.Getenv("DATA_SVC"))

	handler := handler.NewHandler(&storage, redisClient, dataClient)

	router.HandleFunc(constants.CREATE_ROOM, handler.CreateRoom).Methods(http.MethodGet)
	router.HandleFunc(constants.ADD_PLAYER, handler.AddPlayer).Methods(http.MethodGet)
	router.HandleFunc(constants.START_GAME, handler.StartGame).Methods(http.MethodGet)
	router.HandleFunc(constants.GET_MCQ, handler.GetMCQ).Methods(http.MethodGet)
	router.HandleFunc(constants.END_GAME, handler.EndGame).Methods(http.MethodGet)

	router.HandleFunc(constants.ADD_ANSWER, handler.SubmitMCQ).Methods(http.MethodPost)
	router.HandleFunc("/get_result/{room_id}", handler.GetResult).Methods(http.MethodGet)

	http.ListenAndServe(os.Getenv("HTTP_PORT"), router)
}
