package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/mcq_backend/constants"
	"github.com/mcq_backend/models"
	"github.com/mcq_backend/nats"
	"github.com/mcq_backend/storage"
	"github.com/mcq_backend/utils"
)

type Handler struct {
	storage *storage.Storage
	ns      *nats.NatsStruct
}

func NewHandler(st *storage.Storage) *Handler {
	ns, err := nats.NewNats()
	if err != nil {
		log.Fatalln(err)
	}

	return &Handler{
		storage: st,
		ns:      ns,
	}
}

func (h *Handler) GetMCQ(w http.ResponseWriter, r *http.Request) {
	utils.SetHttpHeaders(&w)

	jsonArray, err := utils.ReadJSON()

	if err != nil {
		fmt.Println("Error in ReadJSON", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(jsonArray)
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	utils.SetHttpHeaders(&w)

	id := utils.GenerateUUID()
	reasonCode := h.storage.AddRoom(id)

	if reasonCode == constants.ROOM_CREATED {
		json.NewEncoder(w).Encode(&models.McqAnswer{Room: id})
		w.WriteHeader(http.StatusOK)

		return
	}

	constants.WriteErrorResponse(constants.ROOM_EXIST, "Room already exist", &w)
	w.WriteHeader(http.StatusBadRequest)
}

func (h *Handler) AddPlayer(w http.ResponseWriter, r *http.Request) {
	utils.SetHttpHeaders(&w)

	params := mux.Vars(r)
	roomID := params[constants.ROOM_ID]
	playerID := params[constants.PLAYER_ID]
	isAdmin, _ := strconv.Atoi(params[constants.ADMIN])

	responseCode := h.storage.AddPlayer(roomID, playerID, isAdmin)

	if responseCode == constants.ROOM_NOT_FOUND {
		w.WriteHeader(http.StatusBadRequest)
		constants.WriteErrorResponse(constants.ROOM_NOT_FOUND, "Room not found", &w)

		return
	}

	if responseCode == constants.GAME_STARTED {
		w.WriteHeader(http.StatusBadRequest)
		constants.WriteErrorResponse(constants.GAME_STARTED, "Game has started", &w)

		return
	}

	if responseCode == constants.PLAYER_EXIST {
		w.WriteHeader(http.StatusBadRequest)
		constants.WriteErrorResponse(constants.PLAYER_EXIST, "Player existed", &w)

		return
	}

	// Publising the nats message because of addition of player
	players := []string{}
	currentPlayersInRoom, _ := h.storage.ReturnRoomDetails(roomID)

	for player := range currentPlayersInRoom {
		players = append(players, player)
	}

	h.ns.PlayerJoinMessage(roomID, players)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) SubmitMCQ(w http.ResponseWriter, r *http.Request) {
	utils.SetHttpHeaders(&w)

	var answerSet models.McqAnswer

	err := json.NewDecoder(r.Body).Decode(&answerSet)

	log.Println("Submitted answwer set", answerSet, answerSet.AnswerSet)

	if err != nil {
		fmt.Println("Error in Decode", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	responseCode := h.storage.AddAnswer(answerSet.Room, answerSet.Player, answerSet.AnswerSet)

	if responseCode == constants.ROOM_NOT_FOUND {
		w.WriteHeader(http.StatusBadRequest)
		constants.WriteErrorResponse(constants.ROOM_NOT_FOUND, "Room not found", &w)

		return
	}

	if responseCode == constants.PLAYER_NOT_FOUND {
		w.WriteHeader(http.StatusBadRequest)
		constants.WriteErrorResponse(constants.PLAYER_NOT_FOUND, "Player not found", &w)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) EvaluteResult(roomID string) *models.ScoreCard {
	players, reasonCode := h.storage.ReturnRoomDetails(roomID)

	if reasonCode == constants.ROOM_NOT_FOUND {
		return nil
	}

	if len(players) == 0 {
		return nil
	}

	mcqArray, err := utils.ReadJSON()

	if err != nil {
		fmt.Println("Error in ReadJSON", err)
		return nil
	}

	result := models.ScoreCard{
		ScoreCard: []models.PlayerPerformance{},
	}

	for player, answerSet := range players {
		pc := models.PlayerPerformance{}
		correctAnswer := 0
		if strings.HasPrefix(player, "__") {
			continue
		}

		if len(answerSet) == 1 {
			correctAnswer = -1
		} else {
			for id, answer := range answerSet {
				if strings.HasPrefix(id, "__") {
					continue
				}

				id, _ := strconv.Atoi(id)
				mcqSlice := mcqArray.QuestionSet[id-1]
				mcqCorrect := mcqSlice.Answer

				if answer == mcqCorrect {
					correctAnswer += 1
				}
			}
		}
		pc.PlayerID = player
		pc.CorrectAnswer = correctAnswer

		result.ScoreCard = append(result.ScoreCard, pc)
	}

	return &result
}

func (h *Handler) SSE(w http.ResponseWriter, r *http.Request) {
	utils.SetHttpHeaders(&w)

	params := mux.Vars(r)
	roomID := params[constants.ROOM_ID]

	result := h.EvaluteResult(roomID)

	if result != nil {
		json.NewEncoder(w).Encode(result)
	}
}

func (h *Handler) StartGame(w http.ResponseWriter, r *http.Request) {
	utils.SetHttpHeaders(&w)

	params := mux.Vars(r)
	roomID := params[constants.ROOM_ID]
	endTime, _ := strconv.Atoi(params[constants.ENDGAME_ID])

	reasonCode := h.storage.StartGame(roomID, endTime)
	if reasonCode&constants.ROOM_NOT_FOUND == constants.ROOM_NOT_FOUND {
		log.Fatalf("error in store start game room not found")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err := h.ns.PlayerStartGame(roomID)
	if err != nil {
		log.Fatalf("error in nats start game %+v \n", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	delay := time.Duration(endTime) * time.Millisecond
	log.Println("calling end after", delay)
	//starting counter to end the game at end
	go func() {
		time.Sleep(delay) // Adjust the duration as needed

		// Execute the last two lines
		log.Println(">> sending end_game nats message")
		h.ns.EndGame(roomID)

		time.Sleep(10 * time.Second)
		log.Println(">> deelting data from server")
		h.storage.FlushRoom(roomID)
		h.ns.DeleteTopic(roomID)
	}()
}

func (h *Handler) EndGame(w http.ResponseWriter, r *http.Request) {
	utils.SetHttpHeaders(&w)

	params := mux.Vars(r)

	roomID := params[constants.ROOM_ID]
	endTime := params[constants.TIME]

	log.Println("Calling endGame", roomID, endTime)
	defer log.Println("Exiting endGame")

	delay, err := time.ParseDuration(endTime + "ms")
	if err != nil {
		log.Fatalln("Error parsing endTime:", err)
		return
	}

	go func() {
		// Wait for a specific duration before executing the following lines
		time.Sleep(delay) // Adjust the duration as needed

		// Execute the last two lines
		h.storage.FlushRoom(roomID)
		h.ns.DeleteTopic(roomID)
	}()
}
