package internal

import (
	"encoding/json"
	"net/http"
)

type Status struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func SetAnswer(s Status, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func MessageByMap(status string) map[string]interface{} {
	m := make(map[string]interface{})
	m["status"] = status
	return m
}

func SetJSON(data interface{}, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
