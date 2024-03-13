package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Actor struct {
	Name     string `json:"name"`
	Sex      int    `json:"sex"`
	Birthday string `json:"birthday"`
}

type Film struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	Rating      int     `json:"rating"`
	Actors      []Actor `json:"actors"`
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     int    `json:"role"` // 0 - user, 1 - admin
}

type Storage struct {
	db *sql.DB
}

func NewDB() (Storage, error) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=770948 dbname=vk sslmode=disable")
	if err != nil {
		return Storage{}, err
	}

	err = db.Ping()
	if err != nil {
		return Storage{}, err
	}

	StartSQL(db)

	return Storage{db: db}, nil
}

func StartSQL(db *sql.DB) {
	var DDL = []string{`CREATE TABLE IF NOT EXISTS films (
		film_id serial PRIMARY KEY,
		name VARCHAR(150), 
		description VARCHAR(1000),
		date DATE,
		rating INT CHECK(rating BETWEEN 1 AND 10)
		);`,
		`CREATE TABLE IF NOT EXISTS actors (
			actor_id serial PRIMARY KEY,
			sex  CHAR(1) CHECK(sex='M' OR sex='W'),
			birthday DATE
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			user_id  serial PRIMARY KEY,
			name  VARCHAR(60) UNIQUE,
			password VARCHAR,
			role  INTEGER DEFAULT 0 CHECK(role >=  0 AND role <= 1)
		);`}

	for i := 0; i < len(DDL); i++ {
		pc, errP := db.Prepare(DDL[i]) // prepare statement but do not execute it yet
		if errP != nil {
			log.Fatal("Error preparing ddl: ", errP)
		}
		_, err := pc.Exec()
		if err != nil {
			log.Fatal("Error creating table: ", err)
		}
	}
}
