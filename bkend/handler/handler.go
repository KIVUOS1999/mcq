package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/mcq_backend/constants"
	"github.com/mcq_backend/models"
	"github.com/mcq_backend/nats"
	"github.com/mcq_backend/storage"
	"github.com/mcq_backend/utils"

	dataClient "github.com/bkend-db/pkg/client"
	restServiceModels "github.com/bkend-redis/models"
	redisClient "github.com/bkend-redis/pkg/client"
)

type Handler struct {
	redisService *redisClient.Client
	dataService  *dataClient.Client
	storage      *storage.Storage
	ns           *nats.NatsStruct
}

func NewHandler(st *storage.Storage, redisSvc *redisClient.Client, dataSvc *dataClient.Client) *Handler {
	ns, err := nats.NewNats()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Nats connected")
	return &Handler{
		storage:      st,
		ns:           ns,
		redisService: redisSvc,
		dataService:  dataSvc,
	}
}

func (h *Handler) createQuestionSet(roomID string, params map[string]string, w http.ResponseWriter) {
	questionCount := params[constants.QUESTION_COUNT]

	count, err := strconv.Atoi(questionCount)
	if err != nil {
		count = 10
	}

	// get the questions set from the data layer
	resp, err := h.dataService.GetQuestionSet(count)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	var questionSet restServiceModels.McqArray

	err = json.NewDecoder(resp.Body).Decode(&questionSet)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	log.Printf("Question set: %+v\n", questionSet)

	// add it to the redis server so that all the player in room can use it
	resp, err = h.redisService.AddQuestionSet(roomID, &questionSet)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.Copy(w, resp.Body)
		return
	}

	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(questionSet)
}

func (h *Handler) getQuestionSet(roomID string, w http.ResponseWriter) {
	resp, err := h.redisService.GetQuestionSet(roomID)
	if err != nil {
		log.Println("err:", err.Error())

		return
	}

	var questionSet restServiceModels.McqArray

	err = json.NewDecoder(resp.Body).Decode(&questionSet)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	log.Printf("Get question set; %+v", questionSet)

	json.NewEncoder(w).Encode(questionSet)
}

func (h *Handler) GetMCQ(w http.ResponseWriter, r *http.Request) {
	// admin must me removed from near future as the responsiblity is to only get the question set not create the question set.
	utils.SetHttpHeaders(&w)

	params := mux.Vars(r)
	roomID := params[constants.ROOM_ID]
	admin, err := strconv.ParseBool(params[constants.ADMIN])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	if roomID == "" {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if admin {
		h.createQuestionSet(roomID, params, w)

		return
	}

	h.getQuestionSet(roomID, w)
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	utils.SetHttpHeaders(&w)

	log.Print("CreateRoom - Enter")
	defer log.Print("CreateRoom - Exit")

	id := utils.GenerateUUID()
	resp, err := h.redisService.CreateRoom(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"err": err.Error()})

		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)

		io.Copy(w, resp.Body)
		return
	}

	json.NewEncoder(w).Encode(&models.McqAnswer{Room: id})
}

func (h *Handler) AddPlayer(w http.ResponseWriter, r *http.Request) {
	utils.SetHttpHeaders(&w)

	params := mux.Vars(r)
	roomID := params[constants.ROOM_ID]
	playerID := params[constants.PLAYER_ID]
	isAdmin := params[constants.ADMIN]

	resp, err := h.redisService.AddPlayer(roomID, playerID, isAdmin)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"err": err.Error()})

		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
		return
	}

	// Publising the nats message because of addition of player
	currentPlayersInRoom, err := h.redisService.GetPlayer(roomID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"err": err.Error()})

		return
	}

	if currentPlayersInRoom.StatusCode != http.StatusOK {
		w.WriteHeader(currentPlayersInRoom.StatusCode)
		io.Copy(w, currentPlayersInRoom.Body)

		return
	}

	playersArr := []string{}

	err = json.NewDecoder(currentPlayersInRoom.Body).Decode(&playersArr)

	h.ns.PlayerJoinMessage(roomID, playersArr)
}

func (h *Handler) SubmitMCQ(w http.ResponseWriter, r *http.Request) {
	utils.SetHttpHeaders(&w)

	var answerSet models.McqAnswer

	err := json.NewDecoder(r.Body).Decode(&answerSet)

	log.Println("Submitted answer set", answerSet, answerSet.AnswerSet)

	if err != nil {
		fmt.Println("Error in Decode", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	resp, err := h.redisService.AddAnswers(answerSet.Room, answerSet.Player, answerSet.AnswerSet)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"err": err.Error()})

		return
	}

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) evaluteResult(roomID string) *models.ScoreCard {
	resp, err := h.redisService.GetAnswers(roomID)

	if err != nil {
		log.Println("Error in restService.GetAnswers", err.Error())
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("Status not 200 OK")
		return nil
	}

	// here we are having answers that players submitted.
	// {player: {answer: ["question_id":player_answer]}}
	var roomStruct map[string]restServiceModels.AnswersListStruct
	err = json.NewDecoder(resp.Body).Decode(&roomStruct)
	if err != nil {
		log.Println("error in decoder", err.Error())
		return nil
	}

	log.Printf("%+v - %+v", roomStruct, resp.Body)

	players := roomStruct

	// getting the real answers
	resp, err = h.redisService.GetQuestionSet(roomID)
	if err != nil {
		log.Println("Error in redis getQuestionSet", err.Error())
	}

	var questionSet restServiceModels.McqArray

	err = json.NewDecoder(resp.Body).Decode(&questionSet)
	if err != nil {
		fmt.Println("Error in ReadJSON", err)
		return nil
	}

	correctAnswerList := map[string]int{}
	pcd := models.PlayerPerformanceDetail{
		DetailQusetionSet: map[string]models.PlayerPerformanceDetailQuestionSet{},
	}

	detailedAnswerSet := models.PlayerPerformanceDetailQuestionSet{}
	for _, val := range questionSet.QuestionSet {
		intAnswer, _ := strconv.Atoi(val.Answer)
		correctAnswerList[val.ID] = intAnswer

		detailedAnswerSet.Question = val.Question
		detailedAnswerSet.CorrectOption = intAnswer
		detailedAnswerSet.CorrectAnswer = val.Options[val.Answer]

		pcd.DetailQusetionSet[val.ID] = detailedAnswerSet
	}

	log.Println("Correct answer list", correctAnswerList)

	// comparing what players have submitted and what is the actual answer.
	result := models.ScoreCard{
		ScoreCard: []models.PlayerPerformance{},
		Detail:    []models.PlayerPerformanceDetail{},
	}

	for player, answerSet := range players {
		pc := models.PlayerPerformance{}
		playerPerformanceDetail := pcd
		correctAnswer := 0

		if len(answerSet.Answers) == 1 {
			// if the player has not yet submitted
			correctAnswer = -1
		} else {
			for questionID, submittedAnswer := range answerSet.Answers {
				originalAnswer := correctAnswerList[questionID]
				log.Printf(">> \t %s \t %+v \t %+v \t %+v\n", player, questionID, submittedAnswer, originalAnswer)

				if originalAnswer == submittedAnswer {
					correctAnswer += 1
				}

				filledOptions := playerPerformanceDetail.DetailQusetionSet[questionID]
				filledOptions.SelectedOption = submittedAnswer
				playerPerformanceDetail.DetailQusetionSet[questionID] = filledOptions
			}
		}
		pc.PlayerID = player
		pc.CorrectAnswer = correctAnswer
		playerPerformanceDetail.PlayerID = player

		result.ScoreCard = append(result.ScoreCard, pc)
		result.Detail = append(result.Detail, playerPerformanceDetail)
	}

	return &result
}

func (h *Handler) GetResult(w http.ResponseWriter, r *http.Request) {
	utils.SetHttpHeaders(&w)

	params := mux.Vars(r)
	roomID := params[constants.ROOM_ID]

	result := h.evaluteResult(roomID)

	if result == nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("cannot make the player result")

		return
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("marshal error")

		return
	}

	log.Println(string(resultBytes))

	h.ns.Submission(roomID, resultBytes)
}

func (h *Handler) StartGame(w http.ResponseWriter, r *http.Request) {
	utils.SetHttpHeaders(&w)

	params := mux.Vars(r)
	roomID := params[constants.ROOM_ID]
	endTime, _ := strconv.Atoi(params[constants.ENDGAME_ID])

	resp, err := h.redisService.StartGame(roomID, params[constants.ENDGAME_ID])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"err": err.Error()})

		return
	}

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)

		return
	}

	h.createQuestionSet(roomID, map[string]string{constants.QUESTION_COUNT: "10"}, w)

	err = h.ns.PlayerStartGame(roomID, endTime)
	if err != nil {
		log.Printf("error in nats start game %+v \n", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	// setting up the question for the group

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
		h.redisService.FlushRoom(roomID)
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
		h.redisService.FlushRoom(roomID)
		h.ns.DeleteTopic(roomID)
	}()
}
