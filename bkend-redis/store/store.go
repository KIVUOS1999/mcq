package store

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/bkend-redis/models"
	con "github.com/bkend-redis/pkg/constants"
	rediscache "github.com/bkend-redis/redis-cache"
)

type Store struct {
	redisClient *rediscache.CacheStruct
}

func New(redisClient *rediscache.CacheStruct) *Store {
	return &Store{redisClient: redisClient}
}

func (s *Store) roomDetailsString(roomID string) *string {
	val, err := s.redisClient.GetValue(roomID)
	if err != nil || val == "" {
		return nil
	}

	return &val
}

func (s *Store) PlayerDetails(roomID string) (*models.RoomsStruct, int, error) {
	roomData := s.roomDetailsString(roomID)
	if roomData == nil {
		return nil, con.RoomDoesNotExist, errors.New("room does not exists")
	}

	var roomStruct models.RoomsStruct

	if err := json.Unmarshal([]byte(*roomData), &roomStruct); err != nil {
		return nil, con.NotSpecific, err
	}

	return &roomStruct, con.NotSpecific, nil
}

func (s *Store) AddRoom(roomID string) (int, error) {
	if s.roomDetailsString(roomID) != nil {
		return con.RoomExist, errors.New("room already present")
	}

	var roomStruct models.RoomsStruct

	roomStruct.PlayerList = models.PlayersStruct{
		Player: map[string]models.AnswersListStruct{},
	}
	data, err := json.Marshal(roomStruct)
	if err != nil {
		return con.NotSpecific, err
	}

	return con.NotSpecific, s.redisClient.SetValue(roomID, string(data))
}

func (s *Store) AddPlayer(roomID, playerID string, admin bool) (int, error) {
	roomStruct, identifier, err := s.PlayerDetails(roomID)
	if err != nil {
		return identifier, err
	}

	if roomStruct.HasStarted {
		return con.GameStarted, errors.New("game has started")
	}

	if _, ok := roomStruct.PlayerList.Player[playerID]; ok {
		return con.PlayerExist, errors.New("player already exists")
	}

	if roomStruct.Admin == "" && admin {
		roomStruct.Admin = playerID
	}

	roomStruct.PlayerList.Player[playerID] = models.AnswersListStruct{}

	val, err := json.Marshal(roomStruct)
	if err != nil {
		return con.NotSpecific, err
	}

	return con.NotSpecific, s.redisClient.SetValue(roomID, string(val))
}

func (s *Store) StartGame(roomID string, endTime int) (int, error) {
	currTime := time.Now()
	durationInMills := time.Duration(endTime) * time.Microsecond
	endTimeEpoch := currTime.Add(durationInMills)

	roomStruct, errID, err := s.PlayerDetails(roomID)
	if err != nil {
		return errID, err
	}

	roomStruct.HasStarted = true
	roomStruct.EndTime = endTimeEpoch

	val, err := json.Marshal(roomStruct)
	if err != nil {
		return con.NotSpecific, err
	}

	return con.NotSpecific, s.redisClient.SetValue(roomID, string(val))
}

func (s *Store) ReturnRoomDetails(roomID string) (*models.RoomsStruct, int, error) {
	val, err := s.redisClient.GetValue(roomID)
	if err != nil {
		return nil, con.RoomDoesNotExist, err
	}

	var roomStruct models.RoomsStruct

	err = json.Unmarshal([]byte(val), &roomStruct)
	if err != nil {
		return nil, con.NotSpecific, err
	}

	return &roomStruct, con.NotSpecific, nil
}

func (s *Store) AddPlayerAnswers(roomID, playerID string, answers map[string]int) (int, error) {
	roomStruct, errID, err := s.PlayerDetails(roomID)
	if err != nil {
		return errID, err
	}

	if _, ok := roomStruct.PlayerList.Player[playerID]; !ok {
		return con.PlayerNotExist, errors.New("Player doesnot exists")
	}

	roomStruct.PlayerList.Player[playerID] = models.AnswersListStruct{
		Answers: answers,
	}

	val, err := json.Marshal(roomStruct)
	if err != nil {
		return con.NotSpecific, err
	}

	s.redisClient.SetValue(roomID, string(val))

	return con.NotSpecific, nil
}

func (s *Store) AddQuestionSet(roomID string, questionSet *models.McqArray) (int, error) {
	roomStruct, errID, err := s.PlayerDetails(roomID)
	if err != nil {
		return errID, err
	}

	roomStruct.QuestionList = *questionSet

	val, err := json.Marshal(roomStruct)
	if err != nil {
		return con.NotSpecific, err
	}

	s.redisClient.SetValue(roomID, string(val))

	return con.NotSpecific, nil
}

func (s *Store) GetQuestionSet(roomID string) (*models.McqArray, int, error) {
	roomStruct, errID, err := s.PlayerDetails(roomID)
	if err != nil {
		return nil, errID, err
	}

	questionSet := roomStruct.QuestionList

	return &questionSet, con.NotSpecific, nil
}

func (s *Store) AddRoomQuestion(roomID string) (*models.McqArray, int, error) {
	roomStruct, errID, err := s.PlayerDetails(roomID)
	if err != nil {
		return nil, errID, err
	}

	questionSet := roomStruct.QuestionList

	return &questionSet, con.NotSpecific, nil
}

func (s *Store) GetRoomQuestion(roomID string) (*models.McqArray, int, error) {
	roomStruct, errID, err := s.PlayerDetails(roomID)
	if err != nil {
		return nil, errID, err
	}

	questionSet := roomStruct.QuestionList

	return &questionSet, con.NotSpecific, nil
}

func (s *Store) FlushRoom(roomID string) error {
	err := s.redisClient.DeleteKey(roomID)
	if err != nil {
		return err
	}

	return nil
}
