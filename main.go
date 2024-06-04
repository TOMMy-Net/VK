package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/TOMMy-Net/VK/db"
	"github.com/TOMMy-Net/VK/handlers"
	"github.com/TOMMy-Net/VK/middleware"
	"github.com/TOMMy-Net/VK/services"
	"github.com/joho/godotenv"
)

// @title Blueprint Swagger API
// @version 1.0

// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @license.name MIT

// @BasePath /api/

func main() {
	var errENV = godotenv.Load() // load env
	if errENV != nil {
		log.Fatal(errENV)
	}

	var db, errDB = db.NewDB() // load db
	if errDB != nil {
		log.Fatal(errDB)
	}

	var file, err = os.OpenFile("request.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // file for logging
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var auth = services.NewAuthService(db) // load auth

	var authWare = middleware.TokenAuthWare(db, auth) // auth middleware

	var servH = handlers.Service{Storage: db, Auth: auth} // handlers class
	mux := http.NewServeMux()

	mux.Handle("/api/films", authWare(servH.FilmsInformation()))
	mux.Handle("/api/films/search", authWare(servH.SearchFilmHandler()))
	mux.Handle("/api/actors", authWare(servH.ActorsInformation()))
	mux.Handle("/api/auth", authWare(servH.AuthByUserHandler()))
	//mux.Handle("/swagger", ginSwagger.WrapHandler(swaggerFiles.Handler))
	fmt.Println("Start server")
	log.Fatal(http.ListenAndServe(":8000", services.LoggingHandler(file, mux)))
}
