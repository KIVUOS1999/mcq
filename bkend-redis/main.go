package main

import (
	"net/http"
	"os"

	"github.com/bkend-redis/constants"
	"github.com/bkend-redis/handler"
	rediscache "github.com/bkend-redis/redis-cache"
	"github.com/bkend-redis/store"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("configs/.env")
	if err != nil {
		panic("Error loading .env file")
	}

	redisCache := rediscache.New()
	store := store.New(redisCache)
	handler := handler.New(store)

	router := mux.NewRouter()

	router.HandleFunc("/create_room/room/{"+constants.RoomID+"}", handler.CreateRoom).Methods(http.MethodPost)
	router.HandleFunc("/add_player/room/{"+constants.RoomID+"}/player/{"+constants.PlayerID+"}/admin/{"+constants.Admin+"}", handler.AddPlayer).Methods(http.MethodPost)
	router.HandleFunc("/get_player/room/{"+constants.RoomID+"}", handler.GetPlayers).Methods(http.MethodGet)
	router.HandleFunc("/start_game/room/{"+constants.RoomID+"}/end_time/{"+constants.EndTime+"}", handler.StartGame).Methods(http.MethodPost)

	router.HandleFunc("/question_set/room/{"+constants.RoomID+"}", handler.AddQuestionSet).Methods(http.MethodPost)
	router.HandleFunc("/question_set/room/{"+constants.RoomID+"}", handler.GetQuestionSet).Methods(http.MethodGet)

	router.HandleFunc("/add_answer/room/{"+constants.RoomID+"}/player/{"+constants.PlayerID+"}", handler.AddAnswers).Methods(http.MethodPost)
	router.HandleFunc("/get_answers/room/{"+constants.RoomID+"}", handler.GetAnswers).Methods(http.MethodGet)
	router.HandleFunc("/flush_room/room/{"+constants.RoomID+"}", handler.FlushRoom).Methods(http.MethodDelete)

	http.ListenAndServe(os.Getenv("HTTP_PORT"), router)
}
