package constants

import (
	"encoding/json"
	"net/http"
)

type ErrorReason struct {
	ResponseCode int    `json:"response_code"`
	Reason       string `json:"reason"`
}

func WriteErrorResponse(responseCode int, reason string, w *http.ResponseWriter) {
	err := ErrorReason{
		responseCode,
		reason,
	}

	json.NewEncoder(*w).Encode(err)
}
