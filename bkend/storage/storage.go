package storage

import (
	"log"
	"time"

	"github.com/mcq_backend/constants"
)

/*
{
	room1 : {
		player1 : {
			{
				__admin : 1
				q_id_1 : answer from player,
				q_id_2 : answer from player,
				q_id_3 : answer from player,
				...
			}
		},
		player2 : {
			__admin : 0
		}
		player3 : {

		}
		__started : {__started: 0}
		__endTime : {__endTime: end-epoch}
	},
	room2 : {
		...
	},
	...
}
*/

type Storage struct {
	Rooms map[string]map[string]map[string]int `json:"rooms"` //RoomID -> PlayerID -> Answer
}

func NewStorage() Storage {
	return Storage{
		Rooms: map[string]map[string]map[string]int{},
	}
}

func (s *Storage) roomExists(roomID string) (reasonCode uint16) {
	if _, ok := s.Rooms[roomID]; ok {
		return constants.ROOM_EXIST
	}

	return constants.ROOM_NOT_FOUND
}

func (s *Storage) gameStarted(roomID string) (reasonCode uint16) {
	if s.Rooms[roomID][constants.STORE_START_KEY][constants.STORE_START_KEY] == 1 {
		return constants.GAME_STARTED
	}

	return constants.SUCCESS
}

func (s *Storage) playerExists(roomID, playerID string) (reason uint16) {
	if _, ok := s.Rooms[roomID][playerID]; ok {
		return constants.PLAYER_EXIST
	}

	return constants.PLAYER_NOT_FOUND
}

func (s *Storage) AddRoom(roomID string) (reasonCode uint16) {
	if s.roomExists(roomID)&constants.ROOM_EXIST == constants.ROOM_EXIST {
		return constants.ROOM_EXIST
	}
	s.Rooms[roomID] = map[string]map[string]int{}
	s.Rooms[roomID][constants.STORE_START_KEY] = map[string]int{
		constants.STORE_START_KEY: 0,
	}

	return constants.ROOM_CREATED
}

func (s *Storage) StartGame(roomID string, endTime int) (reasonCode uint16) {
	currTime := time.Now()
	durationInMills := time.Duration(endTime) * time.Microsecond
	endTimeEpoch := currTime.Add(durationInMills)

	log.Println(endTimeEpoch)

	if s.roomExists(roomID)&constants.ROOM_NOT_FOUND == constants.ROOM_NOT_FOUND {
		return constants.ROOM_NOT_FOUND
	}

	s.Rooms[roomID][constants.STORE_START_KEY][constants.STORE_START_KEY] = 1
	s.Rooms[roomID][constants.STORE_ENDTIME_KEY] = map[string]int{
		constants.STORE_ENDTIME_KEY: int(endTimeEpoch.Unix()),
	}

	return constants.ROOM_CREATED
}

func (s *Storage) AddPlayer(roomID, playerID string, admin int) (reasonCode uint16) {
	if s.roomExists(roomID)&constants.ROOM_NOT_FOUND == constants.ROOM_NOT_FOUND {
		return constants.ROOM_NOT_FOUND
	}

	if s.gameStarted(roomID)&constants.GAME_STARTED == constants.GAME_STARTED {
		return constants.GAME_STARTED
	}

	if s.playerExists(roomID, playerID)&constants.PLAYER_EXIST == constants.PLAYER_EXIST {
		return constants.PLAYER_EXIST
	}

	s.Rooms[roomID][playerID] = map[string]int{
		constants.STORE_ADMIN_KEY: admin,
	}

	return constants.PLAYER_CREATED
}

func (s *Storage) AddAnswer(roomID, playerID string, answers map[string]int) (reasonCode uint16) {
	if s.roomExists(roomID)&constants.ROOM_NOT_FOUND == constants.ROOM_NOT_FOUND {
		return constants.ROOM_NOT_FOUND
	}

	if s.playerExists(roomID, playerID)&constants.PLAYER_NOT_FOUND == constants.PLAYER_NOT_FOUND {
		return constants.PLAYER_NOT_FOUND
	}

	s.Rooms[roomID][playerID] = answers

	return constants.ANSWER_ADDED
}

func (s *Storage) ReturnRoomDetails(roomID string) (players map[string]map[string]int, reasonCode uint16) {
	if s.roomExists(roomID)&constants.ROOM_NOT_FOUND == constants.ROOM_NOT_FOUND {
		return nil, constants.ROOM_NOT_FOUND
	}

	updatedRoomVal := map[string]map[string]int{}

	for key, val := range s.Rooms[roomID] {
		log.Println("room details:", key, val)
		if key == constants.STORE_START_KEY {
			continue
		}

		updatedRoomVal[key] = val
	}

	return updatedRoomVal, constants.SUCCESS
}

func (s *Storage) FlushRoom(roomID string) {
	delete(s.Rooms, roomID)
}
