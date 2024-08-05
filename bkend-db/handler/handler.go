package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/bkend-db/store"
	"github.com/gorilla/mux"
)

type handler struct {
	store *store.Store
}

func New(s *store.Store) *handler {
	return &handler{store: s}
}

func (h *handler) GetCountedRecords(res http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	count := params["count"]

	c, err := strconv.Atoi(count)
	if err != nil {
		c = 10
	}

	qs, err := h.store.GetRecords(c)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	log.Println(qs)

	jsonData, _ := json.Marshal(map[string]interface{}{
		"data": qs,
	})

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonData)
}

func (h *handler) GetQuestionByID(res http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	questionID := params["id"]

	log.Println("Fetching question for id", questionID)

	qs, err := h.store.GetQuestionByID(questionID)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	log.Println(qs)

	jsonData, _ := json.Marshal(map[string]interface{}{
		"data": qs,
	})

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonData)
}
