package client

import (
	"net/http"
	"strconv"
)

type Client struct {
	url string
}

func New(url string) *Client {
	return &Client{
		url: url,
	}
}

func (c *Client) GetQuestionSet(count int) (*http.Response, error) {
	endPoint := c.url + "/get_question/" + strconv.Itoa(count)

	resp, err := http.Get(endPoint)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetQuestionByID(id string) (*http.Response, error) {
	endPoint := c.url + "/get_question/id/" + id

	resp, err := http.Get(endPoint)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
