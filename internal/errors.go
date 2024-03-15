package internal

import (
	"encoding/json"
	"net/http"
)

const (
	StatusOK = "Ok"
	StatusError = "Error"
)

const(
	JsonError = "Invalid JSON"
	IDError = "Missing parameter 'id'"
	Invalid = "Invalid id"
)

type Error struct {
	Status   string    `json:"status"`
	ErrorMsg string `json:"error"`
}

func ErrorMessage(m string) Error {
	return Error{ErrorMsg: m}
}

func SetErrorJson(err Error, w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
}
