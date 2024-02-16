package constants

const PATH = "./question.json"

// Params
const (
	ROOM_ID   = "room_id"
	PLAYER_ID = "player_id"
)

// Paths
const (
	GET_MCQ        = "/get_mcq"
	SUBMIT_MCQ     = "/submit_mcq"
	CREATE_ROOM    = "/create_room"
	ADD_PLAYER     = "/add_player/room/{" + ROOM_ID + "}/player/{" + PLAYER_ID + "}"
	ADD_ANSWER     = "/add_answer"
	EVALUTE_RESULT = "/evalute_result/{" + ROOM_ID + "}"
)

// Response Code
const (
	SUCCESS = 1
	FAILURE = 0

	ROOM_EXIST     = 2 << 0
	ROOM_NOT_FOUND = 2 << 1
	ROOM_CREATED   = 2 << 2

	PLAYER_EXIST     = 2 << 3
	PLAYER_NOT_FOUND = 2 << 4
	PLAYER_CREATED   = 2 << 5

	ANSWER_ADDED = 2 << 6
)
