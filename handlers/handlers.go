package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TOMMy-Net/VK/db"
	"github.com/TOMMy-Net/VK/internal"
)

type Service struct {
	Storage *db.Storage
}

func (s Service) ActorsInformation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:

		case http.MethodPost:
			if r.URL.Query().Get("id") == "" {
				s.ActorsSetHandler().ServeHTTP(w, r)
			} else {
				s.ActorsUpdateHandler().ServeHTTP(w, r)
			}
		case http.MethodDelete:
			s.ActorsDeleteHandler().ServeHTTP(w, r)
		}
	}
}

func (s Service) ActorsSetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var actor db.Actor
		err := json.NewDecoder(r.Body).Decode(&actor)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: internal.JsonError}, w)
			return
		}
		err = s.Storage.SetActor(actor)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: err.Error()}, w)
			return
		}
		internal.SetAnswer(internal.Status{Status: internal.StatusOK, Message: fmt.Sprintf("Create actor %s success", actor.Name)}, w)
	}
}

func (s Service) ActorsDeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: internal.IDError}, w)
			return
		}
		idI, err := strconv.Atoi(id)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: internal.Invalid}, w)
			return
		}
		count, err := s.Storage.DeleteActor(idI)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: err.Error()}, w)
			return
		} else if count == 0 {
			internal.SetAnswer(internal.Status{Status: internal.StatusOK, Message: fmt.Sprintf("id: %s not in table", id)}, w)
			return
		}
		internal.SetAnswer(internal.Status{Status: internal.StatusOK, Message: fmt.Sprintf("Delete actor id: %s success", id)}, w)
	}
}

func (s Service) ActorsUpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var actor db.Actor
		id := r.URL.Query().Get("id")

		err := json.NewDecoder(r.Body).Decode(&actor)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: internal.JsonError}, w)
			return
		}
		idI, err := strconv.Atoi(id)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: internal.Invalid}, w)
			return
		}
		actor.ID = idI
		count, err := s.Storage.UpdateActor(actor)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: err.Error()}, w)
			return
		} else if count == 0 {
			internal.SetAnswer(internal.Status{Status: internal.StatusOK, Message: fmt.Sprintf("id: %s not in table", id)}, w)
			return
		}
		internal.SetAnswer(internal.Status{Status: internal.StatusOK, Message: fmt.Sprintf("Update actor id: %s success", id)}, w)
	}
}

func (s Service) FilmsInformation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.FilmSetHandler().ServeHTTP(w, r)
		case http.MethodDelete:
			s.FilmDeleteHandler().ServeHTTP(w, r)
		}
	}
}

func (s Service) FilmSetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var film db.Film
		err := json.NewDecoder(r.Body).Decode(&film)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: internal.JsonError}, w)
			return
		}
		err = db.ValidStructFilm(&film)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: err.Error()}, w)
			return
		}
		err = s.Storage.SetFilm(film)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: err.Error()}, w)
			return
		}
		internal.SetAnswer(internal.Status{Status: internal.StatusOK, Message: fmt.Sprintf("Create film %s success", film.Title)}, w)

	}
}

func (s Service) FilmDeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: internal.IDError}, w)
			return
		}
		idI, err := strconv.Atoi(id)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: internal.Invalid}, w)
			return
		}
		count, err := s.Storage.DeleteFilm(idI)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: err.Error()}, w)
			return
		} else if count == 0 {
			internal.SetAnswer(internal.Status{Status: internal.StatusOK, Message: fmt.Sprintf("id: %s not in table", id)}, w)
			return
		}
		internal.SetAnswer(internal.Status{Status: internal.StatusOK, Message: fmt.Sprintf("Delete film id: %s success", id)}, w)
	}
}
