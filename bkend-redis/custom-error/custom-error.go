package customerror

import (
	"encoding/json"
	"net/http"

	customerr "github.com/bkend-redis/pkg/custom-err"
)

func GenerateError(statusCode int, identifier int, reason string, w http.ResponseWriter) {
	err := customerr.CustomError{
		StatusCode: statusCode,
		Identifier: identifier,
		Reason:     reason,
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(err)
}
