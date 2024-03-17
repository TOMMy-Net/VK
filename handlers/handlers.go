package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TOMMy-Net/VK/db"
	"github.com/TOMMy-Net/VK/internal"
	"github.com/TOMMy-Net/VK/services"
)

type Service struct {
	Storage *db.Storage
	Auth    *services.AuthService
}

func (s Service) ActorsInformation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.ActorsGetHandler().ServeHTTP(w, r)
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


// @Summary Getting a list of actors with their films
// @Produce json
// @Success 200 {object} db.Actors
// @Router /api/actors [get]
func (s Service) ActorsGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actors, err := s.Storage.GetActors()
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: err.Error()}, w)
			return
		}
		m := internal.MessageByMap(internal.StatusOK)
		m["actors"] = actors
		internal.SetJSON(m, w)
	}
}



// @Summary      Set actor inforamation
// @Accept       json
// @Produce      json
// @Param        name   sex      birthday
// @Success      200  {object}  model.Actor
// @Failure      400  {object}  httputil.HTTPError
// @Router       /api/actors [post]
func (s Service) ActorsSetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var actor db.Actor
		err := json.NewDecoder(r.Body).Decode(&actor)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: internal.JsonError}, w)
			return
		}
		err = db.ValidStruct(&actor)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: err.Error()}, w)
			return
		}
		err = s.Storage.SetActor(actor)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: err.Error()}, w)
			return
		}
		internal.SetAnswer(internal.Status{Status: internal.StatusOK, Message: fmt.Sprintf("create actor %s success", actor.Name)}, w)
	}
}


// @Summary      delete actor
// @Tags         actors
// @Accept       json
// @Produce      json
// @Param        id     "name search by q"  Format(email)
// @Success      200  {array}   db.Actors
// @Failure      400  {object}  httputil.HTTPError
// @Router       /api/actors [delete]
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
		internal.SetAnswer(internal.Status{Status: internal.StatusOK, Message: fmt.Sprintf("delete actor id: %s success", id)}, w)
	}
}


// @Summary      List accounts
// @Description  get accounts
// @Tags         actors
// @Accept       json
// @Produce      json
// @Param        q    query     string  false  "name search by q"  Format(email)
// @Success      200  {array}   db.Actors
// @Failure      400  {object}  httputil.HTTPError
// @Router       /api/actors?id= [post]
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
		internal.SetAnswer(internal.Status{Status: internal.StatusOK, Message: fmt.Sprintf("update actor id: %s success", id)}, w)
	}
}


func (s Service) FilmsInformation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.GetFilmsBySortHandler().ServeHTTP(w, r)

		case http.MethodPost:
			if r.URL.Query().Get("id") == "" {
				s.FilmSetHandler().ServeHTTP(w, r)
			} else {
				s.FilmUpdateHandler().ServeHTTP(w, r)
			}
		case http.MethodDelete:
			s.FilmDeleteHandler().ServeHTTP(w, r)
		}
	}
}


// @Summary     Set film
// @Tags         films
// @Accept       json
// @Produce      json
// @Param        title    description    date (0000-00-00) rating actors(string)
// @Success      200  {array}   db.Film
// @Failure      400  {object}  httputil.HTTPError
// @Router       /api/films [post]
func (s Service) FilmSetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var film db.Film
		err := json.NewDecoder(r.Body).Decode(&film)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: internal.JsonError}, w)
			return
		}
		err = db.ValidStruct(&film)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: err.Error()}, w)
			return
		}
		err = s.Storage.SetFilm(film)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: err.Error()}, w)
			return
		}
		internal.SetAnswer(internal.Status{Status: internal.StatusOK, Message: fmt.Sprintf("create film %s success", film.Title)}, w)

	}
}

// @Summary     delete film
// @Tags         films
// @Accept       json
// @Produce      json
// @Param        title    description    date (0000-00-00) rating actors(string)
// @Success      200  {json}   db.Film
// @Failure      400  {object}  httputil.HTTPError
// @Router       /api/films?id= [delete]
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
		internal.SetAnswer(internal.Status{Status: internal.StatusOK, Message: fmt.Sprintf("delete film id: %s success", id)}, w)
	}
}

// @Summary     Update film
// @Tags         films
// @Accept       json
// @Produce      json
// @Param        title    description    date (0000-00-00) rating actors(string)
// @Success      200  {json}   db.Film
// @Failure      400  {object}  httputil.HTTPError
// @Router       /api/films?id= [post]
func (s Service) FilmUpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var film db.Film
		id := r.URL.Query().Get("id")

		err := json.NewDecoder(r.Body).Decode(&film)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: internal.JsonError}, w)
			return
		}
		idI, err := strconv.Atoi(id)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: internal.Invalid}, w)
			return
		}
		film.ID = idI
		count, err := s.Storage.UpdateFilm(film)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: err.Error()}, w)
			return
		} else if count == 0 {
			internal.SetAnswer(internal.Status{Status: internal.StatusOK, Message: fmt.Sprintf("id: %s not in table", id)}, w)
			return
		}
		internal.SetAnswer(internal.Status{Status: internal.StatusOK, Message: fmt.Sprintf("update film id: %s success", id)}, w)
	}
}

// @Summary     Sort film
// @Tags         films
// @Accept       json
// @Produce      json
// @Param        type(asc, desc) field
// @Success      200  {array}   db.Film
// @Failure      400  {object}  httputil.HTTPError
// @Router       /api/films?type=&field= [get]
func (s Service) GetFilmsBySortHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sortType := r.URL.Query().Get("type")
		field := r.URL.Query().Get("field")
		if sortType == "" || sortType != db.SortMethodUp {
			sortType = db.SortMethodDown
		}
		if field == "" {
			field = "rating"
		}
		films, err := s.Storage.GetFilmsBySort(field, sortType)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: err.Error()}, w)
			return
		}
		m := internal.MessageByMap(internal.StatusOK)
		m["films"] = films
		internal.SetJSON(m, w)
	}
}

// @Summary     Search film
// @Tags         films
// @Accept       json
// @Produce      json
// @Param        tag
// @Success      200  {array}   db.Film
// @Failure      400  {object}  httputil.HTTPError
// @Router       /api/films/search?tag= [get]
func (s Service) SearchFilmHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := r.URL.Query().Get("tag")

		if tag == "" {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: internal.TagError}, w)
			return
		}
		film, err := s.Storage.GetFilmsBySearch(tag)

		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: err.Error()}, w)
			return
		}
		m := internal.MessageByMap(internal.StatusOK)
		m["films"] = film
		internal.SetJSON(m, w)
	}
}
