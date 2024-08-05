package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/bkend-redis/models"
)

type Client struct {
	url string
}

func New(url string) *Client {
	return &Client{
		url: url,
	}
}

func (c *Client) CreateRoom(roomID string) (*http.Response, error) {
	endPoint := c.url + "/create_room/room/" + roomID
	resp, err := http.Post(endPoint, "application/json", nil)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) AddPlayer(roomID, playerID, isAdmin string) (*http.Response, error) {
	endPoint := c.url + "/add_player/room/" + roomID + "/player/" + playerID + "/admin/" + isAdmin
	resp, err := http.Post(endPoint, "application/json", nil)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetPlayer(roomID string) (*http.Response, error) {
	endPoint := c.url + "/get_player/room/" + roomID
	resp, err := http.Get(endPoint)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) StartGame(roomID, endTime string) (*http.Response, error) {
	endPoint := c.url + "/start_game/room/" + roomID + "/end_time/" + endTime
	resp, err := http.Post(endPoint, "application/json", nil)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) AddAnswers(roomID, playerID string, answers map[string]int) (*http.Response, error) {
	endPoint := c.url + "/add_answer/room/" + roomID + "/player/" + playerID
	body, _ := json.Marshal(answers)
	resp, err := http.Post(endPoint, "application/json", bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetAnswers(roomID string) (*http.Response, error) {
	endPoint := c.url + "/get_answers/room/" + roomID
	resp, err := http.Get(endPoint)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) FlushRoom(roomID string) (*http.Response, error) {
	cli := &http.Client{}
	endPoint := c.url + "/flush_room/room/" + roomID

	req, err := http.NewRequest(http.MethodDelete, endPoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) AddQuestionSet(roomID string, questionSet *models.McqArray) (*http.Response, error) {
	cli := &http.Client{}
	endPoint := c.url + "/question_set/room/" + roomID

	b, _ := json.Marshal(questionSet)

	req, err := http.NewRequest(http.MethodPost, endPoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetQuestionSet(roomID string) (*http.Response, error) {
	cli := &http.Client{}
	endPoint := c.url + "/question_set/room/" + roomID

	req, err := http.NewRequest(http.MethodGet, endPoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
