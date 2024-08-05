package models

type McqStruct struct {
	ID       string            `json:"id"`
	Question string            `json:"question,omitempty"`
	Options  map[string]string `json:"options,omitempty"`
	Answer   string            `json:"answer,omitempty"`
}

type McqArray struct {
	QuestionSet []McqStruct `json:"data,omitempty"`
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

type PlayerPerformanceDetailQuestionSet struct {
	Question       string `json:"question"`
	SelectedOption int    `json:"selected_option"`
	SelectedAnswer string `json:"selected_answer"`
	CorrectOption  int    `json:"correct_option"`
	CorrectAnswer  string `json:"correct_answer"`
}

type PlayerPerformanceDetail struct {
	PlayerID          string                                        `json:"player_id"`
	DetailQusetionSet map[string]PlayerPerformanceDetailQuestionSet `json:"detailed_answer_set"`
}

type ScoreCard struct {
	ScoreCard []PlayerPerformance       `json:"score_card"`
	Detail    []PlayerPerformanceDetail `json:"detail"`
}
