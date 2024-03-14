package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Actor struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Sex      string `json:"sex"`
	Birthday string `json:"birthday"`
}

type Film struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	Rating      int     `json:"rating"`
	Actors      []Actor `json:"actors"`
}

type User struct {
	ID       int    `json:"id"`
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
		name VARCHAR(150) NOT NULL, 
		description VARCHAR(1000) DEFAULT NULL,
		date DATE DEFAULT NULL,
		rating INT CHECK(rating BETWEEN 1 AND 10)
		);`,
		`CREATE TABLE IF NOT EXISTS actors (
			actor_id serial PRIMARY KEY,
			name VARCHAR UNIQUE NOT NULL,
			sex  CHAR CHECK(sex='M' OR sex='W' OR sex=NULL) DEFAULT NULL,
			birthday DATE DEFAULT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			user_id  serial PRIMARY KEY,
			name  VARCHAR(60) UNIQUE DEFAULT NULL,
			password VARCHAR NOT NULL,
			role  INTEGER DEFAULT 0 CHECK(role >=  0 AND role <= 1)
		);`}

	for i := 0; i < len(DDL); i++ {
		pc, errP := db.Prepare(DDL[i])
		if errP != nil {
			log.Fatal("Error preparing ddl: ", errP)
		}
		_, err := pc.Exec()
		if err != nil {
			log.Fatal("Error creating table: ", err)
		}
	}
}

func (s Storage) SetActor(a Actor) error {
	pc, errP := s.db.Prepare(`INSERT INTO actors (name, sex, birthday) VALUES ($1, $2, $3)`)
	if errP != nil {
		return errP
	}
	_, err := pc.Exec(a.Name, a.Sex, a.Birthday)
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) GetActor(name string) {

}

func (s Storage) DeleteActor(id int) (int, error) {
	pc, errP := s.db.Prepare(`DELETE FROM actors WHERE actor_id = $1`)
	if errP != nil {
		return 0, errP
	}
	r, err := pc.Exec(id)

	if err != nil {
		return 0, err
	}
	count, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (s Storage) UpdateActor(a Actor) (int, error) {
	pc, errP := s.db.Prepare(`UPDATE actors
							SET name =  $1,
							sex = $2,
							birthday = $3
							WHERE actor_id = $4`)
	if errP != nil {
		return 0, errP
	}
	r, err := pc.Exec(a.Name, a.Sex, a.Birthday, a.ID)

	if err != nil {
		return 0, err
	}
	count, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
