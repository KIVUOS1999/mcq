package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/bkend-redis/constants"
	customerror "github.com/bkend-redis/custom-error"
	"github.com/bkend-redis/models"
	consts "github.com/bkend-redis/pkg/constants"
	"github.com/bkend-redis/store"
	"github.com/gorilla/mux"
)

type HandlerStruct struct {
	store *store.Store
}

func New(store *store.Store) HandlerInterface {
	return &HandlerStruct{
		store: store,
	}
}

type HandlerInterface interface {
	CreateRoom(w http.ResponseWriter, r *http.Request)
	AddPlayer(w http.ResponseWriter, r *http.Request)
	GetPlayers(w http.ResponseWriter, r *http.Request)
	StartGame(w http.ResponseWriter, r *http.Request)
	AddAnswers(w http.ResponseWriter, r *http.Request)
	GetAnswers(w http.ResponseWriter, r *http.Request)
	FlushRoom(w http.ResponseWriter, r *http.Request)
	AddQuestionSet(w http.ResponseWriter, r *http.Request)
	GetQuestionSet(w http.ResponseWriter, r *http.Request)
}

func (h *HandlerStruct) CreateRoom(w http.ResponseWriter, r *http.Request) {
	SetHttpHeaders(&w)

	log.Print("CreateRoom - Enter")
	defer log.Print("CreateRoom - Exit")

	roomID := mux.Vars(r)[constants.RoomID]
	if roomID == "" {
		customerror.GenerateError(http.StatusBadRequest, 0, "RoomId is blank", w)
		return
	}

	errCode, err := h.store.AddRoom(roomID)
	if err != nil {
		customerror.GenerateError(http.StatusBadRequest, errCode, err.Error(), w)
		return
	}
}

func (h *HandlerStruct) AddPlayer(w http.ResponseWriter, r *http.Request) {
	SetHttpHeaders(&w)
	roomID := mux.Vars(r)[constants.RoomID]
	playerID := mux.Vars(r)[constants.PlayerID]
	admin := mux.Vars(r)[constants.Admin]

	if roomID == "" || playerID == "" {
		customerror.GenerateError(http.StatusBadRequest, consts.NotSpecific, "RoomId is blank", w)
		return
	}

	isAdmin := false
	if admin == "true" {
		isAdmin = true
	}

	errCode, err := h.store.AddPlayer(roomID, playerID, isAdmin)
	if err != nil {
		customerror.GenerateError(http.StatusBadRequest, errCode, err.Error(), w)
		return
	}

	return
}

func (h *HandlerStruct) GetPlayers(w http.ResponseWriter, r *http.Request) {
	SetHttpHeaders(&w)
	roomID := mux.Vars(r)[constants.RoomID]

	if roomID == "" {
		customerror.GenerateError(http.StatusBadRequest, consts.NotSpecific, "RoomId is blank", w)
		return
	}

	roomDetails, errCode, err := h.store.ReturnRoomDetails(roomID)
	if err != nil {
		customerror.GenerateError(http.StatusBadRequest, errCode, err.Error(), w)
		return
	}

	playerArr := []string{}
	for player := range roomDetails.PlayerList.Player {
		playerArr = append(playerArr, player)
	}

	json.NewEncoder(w).Encode(playerArr)
}

func (h *HandlerStruct) StartGame(w http.ResponseWriter, r *http.Request) {
	SetHttpHeaders(&w)
	roomID := mux.Vars(r)[constants.RoomID]
	endTimeVar := mux.Vars(r)[constants.EndTime]

	if roomID == "" {
		customerror.GenerateError(http.StatusBadRequest, consts.NotSpecific, "RoomId is blank", w)
		return
	}

	endTime, err := strconv.Atoi(endTimeVar)
	if err != nil {
		customerror.GenerateError(http.StatusInternalServerError, consts.NotSpecific, err.Error(), w)
		return
	}

	errCode, err := h.store.StartGame(roomID, endTime)
	if err != nil {
		customerror.GenerateError(http.StatusBadRequest, errCode, err.Error(), w)
		return
	}

	return
}

func (h *HandlerStruct) AddAnswers(w http.ResponseWriter, r *http.Request) {
	SetHttpHeaders(&w)
	roomID := mux.Vars(r)[constants.RoomID]
	playerID := mux.Vars(r)[constants.PlayerID]

	if roomID == "" || playerID == "" {
		customerror.GenerateError(http.StatusBadRequest, consts.NotSpecific, "RoomId or PlayerId is blank", w)
		return
	}

	var answerSet map[string]int

	err := json.NewDecoder(r.Body).Decode(&answerSet)
	if err != nil {
		customerror.GenerateError(http.StatusInternalServerError, consts.NotSpecific, err.Error(), w)
		return
	}

	errCode, err := h.store.AddPlayerAnswers(roomID, playerID, answerSet)
	if err != nil {
		customerror.GenerateError(http.StatusBadRequest, errCode, err.Error(), w)
		return
	}

	return
}

func (h *HandlerStruct) GetAnswers(w http.ResponseWriter, r *http.Request) {
	SetHttpHeaders(&w)
	roomID := mux.Vars(r)[constants.RoomID]

	if roomID == "" {
		customerror.GenerateError(http.StatusBadRequest, consts.NotSpecific, "RoomID or playerID is blank", w)
		return
	}

	roomStruct, errCode, err := h.store.ReturnRoomDetails(roomID)
	if err != nil {
		customerror.GenerateError(http.StatusBadRequest, errCode, err.Error(), w)
		return
	}

	playersMap := roomStruct.PlayerList.Player
	log.Printf("players %+v", playersMap)
	json.NewEncoder(w).Encode(playersMap)
}

func (h *HandlerStruct) FlushRoom(w http.ResponseWriter, r *http.Request) {
	SetHttpHeaders(&w)
	roomID := mux.Vars(r)[constants.RoomID]
	if roomID == "" {
		customerror.GenerateError(http.StatusBadRequest, consts.NotSpecific, "RoomID is blank", w)
		return
	}

	err := h.store.FlushRoom(roomID)
	if err != nil {
		customerror.GenerateError(http.StatusBadRequest, consts.NotSpecific, err.Error(), w)
		return
	}

	return
}

func (h *HandlerStruct) AddQuestionSet(w http.ResponseWriter, r *http.Request) {
	SetHttpHeaders(&w)
	roomID := mux.Vars(r)[constants.RoomID]
	if roomID == "" {
		customerror.GenerateError(http.StatusBadRequest, consts.NotSpecific, "RoomID is blank", w)
		return
	}

	var mcqArray models.McqArray

	err := json.NewDecoder(r.Body).Decode(&mcqArray)
	if err != nil {
		customerror.GenerateError(http.StatusBadRequest, consts.NotSpecific, "decoder error"+err.Error(), w)
		return
	}

	log.Println("Question set", roomID, mcqArray)
	h.store.AddQuestionSet(roomID, &mcqArray)

	return
}

func (h *HandlerStruct) GetQuestionSet(w http.ResponseWriter, r *http.Request) {
	SetHttpHeaders(&w)
	roomID := mux.Vars(r)[constants.RoomID]
	if roomID == "" {
		customerror.GenerateError(http.StatusBadRequest, consts.NotSpecific, "RoomID is blank", w)
		return
	}

	mcqArray, errCode, err := h.store.GetQuestionSet(roomID)
	if err != nil {
		customerror.GenerateError(http.StatusBadRequest, errCode, err.Error(), w)
		return
	}

	log.Println("Question set", mcqArray)

	json.NewEncoder(w).Encode(mcqArray)
}

func SetHttpHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	(*w).Header().Set("Content-Type", "application/json")
}
