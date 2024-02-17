package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mcq_backend/constants"
	"github.com/mcq_backend/models"
	"github.com/mcq_backend/storage"
	"github.com/mcq_backend/utils"
)

type Handler struct {
	storage *storage.Storage
}

func NewHandler(st *storage.Storage) *Handler {
	return &Handler{
		storage: st,
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

	responseCode := h.storage.AddPlayer(roomID, playerID)

	if responseCode == constants.ROOM_NOT_FOUND {
		w.WriteHeader(http.StatusBadRequest)
		constants.WriteErrorResponse(constants.ROOM_NOT_FOUND, "Room not found", &w)
	}

	if responseCode == constants.PLAYER_EXIST {
		w.WriteHeader(http.StatusBadRequest)
		constants.WriteErrorResponse(constants.PLAYER_EXIST, "Player existed", &w)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) SubmitMCQ(w http.ResponseWriter, r *http.Request) {
	utils.SetHttpHeaders(&w)

	var answerSet models.McqAnswer

	err := json.NewDecoder(r.Body).Decode(&answerSet)

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

		if len(answerSet) == 0 {
			correctAnswer = -1
		} else {
			for id, answer := range answerSet {
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

	fmt.Println("Sending SSE")

	params := mux.Vars(r)
	roomID := params[constants.ROOM_ID]

	result := h.EvaluteResult(roomID)

	if result != nil {
		json.NewEncoder(w).Encode(result)
	}
}
