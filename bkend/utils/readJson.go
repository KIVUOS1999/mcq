package utils

import (
	"encoding/json"
	"io"
	"os"

	"github.com/mcq_backend/constants"
	"github.com/mcq_backend/models"
)

func ReadJSON() (*models.McqArray, error) {
	jsonFile, err := os.Open(constants.PATH)
	defer jsonFile.Close()

	if err != nil {
		return nil, err
	}

	jsonContent, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var jsonArray models.McqArray

	err = json.Unmarshal(jsonContent, &jsonArray)
	if err != nil {
		return nil, err
	}

	return &jsonArray, nil
}
