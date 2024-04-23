package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mcq_backend/constants"
	"github.com/mcq_backend/handler"
	"github.com/mcq_backend/storage"
)

func main() {
	router := mux.NewRouter()

	storage := storage.NewStorage()
	handler := handler.NewHandler(&storage)

	router.HandleFunc(constants.GET_MCQ, handler.GetMCQ).Methods("GET")
	router.HandleFunc(constants.ADD_ANSWER, handler.SubmitMCQ).Methods("POST")

	router.HandleFunc(constants.CREATE_ROOM, handler.CreateRoom).Methods("GET")
	router.HandleFunc(constants.ADD_PLAYER, handler.AddPlayer).Methods("GET")
	router.HandleFunc(constants.START_GAME, handler.StartGame).Methods("GET")
	router.HandleFunc(constants.END_GAME, handler.EndGame).Methods("GET")

	router.HandleFunc("/get_result/{room_id}", handler.GetResult)

	http.ListenAndServe(":8000", router)
}
