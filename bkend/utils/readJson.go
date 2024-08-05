package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mcq_backend/models"
)

func ReadJSON() (*models.McqArray, error) {
	url := "http://localhost:8004/get_question/10"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var jsonArray models.McqArray

	err = json.NewDecoder(resp.Body).Decode(&jsonArray)
	if err != nil {
		return nil, err
	}

	fmt.Println(jsonArray)

	return &jsonArray, nil
}
