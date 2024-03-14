package internal

import (
	"encoding/json"
	"net/http"
)

type Status struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func SetAnswer(s Status, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}
