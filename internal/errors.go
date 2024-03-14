package internal

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Status   int    `json:"status"`
	ErrorMsg string `json:"error"`
}

func ErrorMessage(m string) Error {
	return Error{ErrorMsg: m}
}

func SetErrorJson(err Error, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
}
