package main

import (
	"log"
	"net/http"

	"github.com/TOMMy-Net/VK/db"
	"github.com/TOMMy-Net/VK/handlers"
	"github.com/joho/godotenv"
)

func main()  {
	var errENV = godotenv.Load() // load env
	if errENV != nil {
		log.Fatal(errENV)
	}

	db, errDB := db.NewDB() // load db
	if  errDB != nil {
		log.Fatal(errDB)
	}

	var servH = handlers.Service{Storage: &db} 
	mux := http.NewServeMux()

	mux.HandleFunc("/api/films", servH.FilmsInformation())
	mux.HandleFunc("/api/actors", servH.ActorsInformation())
	log.Fatal(http.ListenAndServe(":8000", mux))
}