package main

import (
	"log"
	"net/http"

	"github.com/TOMMy-Net/VK/db"
	"github.com/TOMMy-Net/VK/handlers"
	"github.com/TOMMy-Net/VK/middleware"
	"github.com/TOMMy-Net/VK/services"
	"github.com/joho/godotenv"
)

func main() {
	var errENV = godotenv.Load() // load env
	if errENV != nil {
		log.Fatal(errENV)
	}

	var db, errDB = db.NewDB() // load db
	if errDB != nil {
		log.Fatal(errDB)
	}
	var auth = services.NewAuthService(db) // load auth

	var authWare = middleware.TokenAuthWare(db, auth) // auth middleware

	var servH = handlers.Service{Storage: db, Auth: auth} // handlers class
	mux := http.NewServeMux()

	mux.Handle("/api/films", authWare(servH.FilmsInformation()))
	mux.Handle("/api/films/search", authWare(servH.SearchFilmHandler()))
	mux.Handle("/api/actors", authWare(servH.ActorsInformation()))
	mux.Handle("/api/auth", authWare(servH.AuthByUserHandler()))

	log.Fatal(http.ListenAndServe(":8000", mux))
}
