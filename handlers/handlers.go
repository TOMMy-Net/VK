package handlers


import (
  "net/http"
  "github.com/TOMMy-Net/VK/db"
)

type Service struct {
  Storage     *db.Storage
}


func (s Service) ActorsInformation() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:

    }
  }
}

func (s Service) FilmsInformation() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    
  }
}