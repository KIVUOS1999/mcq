package storage

import (
	"github.com/mcq_backend/constants"
)

/*
{
	room1 : {
		player1 : {
			{
				q_id_1 : answer from player,
				q_id_2 : answer from player,
				q_id_3 : answer from player,
				...
			}
		},
		player2 : {

		}
		player3 : {

		}
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

func (s *Storage) roomExists(roomID string) (reasonCode uint8) {
	if _, ok := s.Rooms[roomID]; ok {
		return constants.ROOM_EXIST
	}

	return constants.ROOM_NOT_FOUND
}

func (s *Storage) playerExists(roomID, playerID string) (reason uint8) {
	if _, ok := s.Rooms[roomID][playerID]; ok {
		return constants.PLAYER_EXIST
	}

	return constants.PLAYER_NOT_FOUND
}

func (s *Storage) AddRoom(roomID string) (reasonCode uint8) {
	if s.roomExists(roomID)&constants.ROOM_EXIST == constants.ROOM_EXIST {
		return constants.ROOM_EXIST
	}
	s.Rooms[roomID] = map[string]map[string]int{}

	return constants.ROOM_CREATED
}

func (s *Storage) AddPlayer(roomID, playerID string) (reasonCode uint8) {
	if s.roomExists(roomID)&constants.ROOM_NOT_FOUND == constants.ROOM_NOT_FOUND {
		return constants.ROOM_NOT_FOUND
	}

	if s.playerExists(roomID, playerID)&constants.PLAYER_EXIST == constants.PLAYER_EXIST {
		return constants.PLAYER_EXIST
	}

	s.Rooms[roomID][playerID] = map[string]int{}

	return constants.PLAYER_CREATED
}

func (s *Storage) AddAnswer(roomID, playerID string, answers map[string]int) (reasonCode uint8) {
	if s.roomExists(roomID)&constants.ROOM_NOT_FOUND == constants.ROOM_NOT_FOUND {
		return constants.ROOM_NOT_FOUND
	}

	if s.playerExists(roomID, playerID)&constants.PLAYER_NOT_FOUND == constants.PLAYER_NOT_FOUND {
		return constants.PLAYER_NOT_FOUND
	}

	if len(s.Rooms[roomID][playerID]) == 0 {
		s.Rooms[roomID][playerID] = answers
	}

	return constants.ANSWER_ADDED
}

func (s *Storage) ReturnRoomDetails(roomID string) (players map[string]map[string]int, reasonCode uint8) {
	if s.roomExists(roomID)&constants.ROOM_NOT_FOUND == constants.ROOM_NOT_FOUND {
		return nil, constants.ROOM_NOT_FOUND
	}

	return s.Rooms[roomID], constants.SUCCESS
}

func (s *Storage) FlushRoom(roomID string) {
	delete(s.Rooms, roomID)
}
