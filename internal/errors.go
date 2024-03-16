package internal

import (
	"encoding/json"
	"net/http"
)

const (
	StatusOK = "OK"
	StatusError = "Error"
)

const(
	JsonError = "invalid JSON"
	IDError = "missing parameter 'id'"
	TagError = "missing parameter 'tag'"
	Invalid = "invalid id"
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
