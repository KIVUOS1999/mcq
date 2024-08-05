package models

import "time"

/*
{
  "player_list": {
    "Player": {
      "player1": {
        "Answers": {
          "answer1": 3,
          "answer2": 1,
          "answer3": 2
        }
      },
      "player2": {
        "Answers": {
          "answer1": 2,
          "answer2": 3,
          "answer3": 1
        }
      }
    }
  },
  "admin": "admin_user",
  "has_started": true,
  "end_time": "2024-07-22T12:00:00Z"
}
*/

type AnswersListStruct struct {
	Answers map[string]int
}

type PlayersStruct struct {
	Player map[string]AnswersListStruct
}

type McqStruct struct {
	ID       string            `json:"id"`
	Question string            `json:"question,omitempty"`
	Options  map[string]string `json:"options,omitempty"`
	Answer   string            `json:"answer,omitempty"`
}

type McqArray struct {
	QuestionSet []McqStruct `json:"data,omitempty"`
}

type RoomsStruct struct {
	PlayerList   PlayersStruct `json:"player_list"`
	QuestionList McqArray      `json:"question_list"`
	Admin        string        `json:"admin"`
	HasStarted   bool          `json:"has_started"`
	EndTime      time.Time     `json:"end_time"`
}

type QuestionIDList struct {
	QuestionIDs []string `json:"question_ids"`
}
