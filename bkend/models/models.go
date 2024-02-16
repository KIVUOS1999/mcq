package models

type McqStruct struct {
	ID       int               `json:"id"`
	Question string            `json:"question,omitempty"`
	Options  map[string]string `json:"options,omitempty"`
	Answer   int               `json:"answer,omitempty"`
}

type McqArray struct {
	QuestionSet []McqStruct `json:"question_set,omitempty"`
}

type McqAnswer struct {
	Room      string         `json:"room_id,omitempty"`
	Player    string         `json:"player_id,omitempty"`
	AnswerSet map[string]int `json:"answer_set,omitempty"`
}

type PlayerPerformance struct {
	PlayerID      string `json:"player_id"`
	CorrectAnswer int    `json:"correct_answer"`
}

type ScoreCard struct {
	ScoreCard []PlayerPerformance `json:"score_card"`
}
