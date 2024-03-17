package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TOMMy-Net/VK/db"
	"github.com/TOMMy-Net/VK/internal"
)

func (s Service) AuthByUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.GetTokenByUserHandler().ServeHTTP(w, r)
		case http.MethodPost:
			s.CreateUserHandler().ServeHTTP(w, r)
		}
	}
}

// @Summary     Get token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        name    password
// @Success      200  {json}   db.User
// @Failure      400  {object}  httputil.HTTPError
// @Router       /api/auth [get]
func (s Service) GetTokenByUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user db.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: internal.JsonError}, w)
			return
		}
		token, err := s.Auth.SignInUser(user.Name, user.Password)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: err.Error()}, w)
			return
		}
		m := internal.MessageByMap(internal.StatusOK)
		m["token"] = token
		internal.SetJSON(m, w)
	}
}

// @Summary     Create user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        name    password
// @Success      200  {json}   db.User
// @Failure      400  {object}  httputil.HTTPError
// @Router       /api/auth [post]
func (s Service) CreateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user db.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: internal.JsonError}, w)
			return
		}
		err = s.Storage.CreateUser(user)
		if err != nil {
			internal.SetErrorJson(internal.Error{Status: internal.StatusError, ErrorMsg: err.Error()}, w)
			return
		}
		internal.SetAnswer(internal.Status{Status: internal.StatusOK, Message: fmt.Sprintf("create user %s success", user.Name)}, w)
	}
}
