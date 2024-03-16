package db

import (
	"database/sql"

	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var (
	ErrAlreadyIn error = fmt.Errorf("this film is already in the database")
)

const (
	SortMethodUp   = "asc"
	SortMethodDown = "desc"
)


type Actor struct {
	ID       int    `json:"id"`
	Name     string `json:"name" validate:"required"`
	Sex      string `json:"sex"`
	Birthday string `json:"birthday" validate:"required"`
	Films    []Film `json:"films"`
}

type Film struct {
	ID          int    `json:"id"`
	Title       string `json:"title" validate:"required,min=1,max=150"`
	Description string `json:"description" validate:"max=1000"`
	Date        string `json:"date"  validate:"required"`
	Rating      int    `json:"rating"`
	Actors      string `json:"actors"`
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
		film_id SERIAL PRIMARY KEY,
		title VARCHAR(150) UNIQUE NOT NULL, 
		description VARCHAR(1000),
		date DATE,
		rating INT CHECK(rating BETWEEN 1 AND 10),
		actors_list VARCHAR
		);`,
		`CREATE TABLE IF NOT EXISTS actors (
			actor_id SERIAL PRIMARY KEY,
			name VARCHAR UNIQUE NOT NULL,
			sex  CHAR CHECK(sex='M' OR sex='W'),
			birthday DATE 
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			user_id  SERIAL PRIMARY KEY,
			name  VARCHAR(60) UNIQUE NOT NULL,
			password VARCHAR NOT NULL,
			role  INTEGER DEFAULT 0 CHECK(role >= 0 AND role <= 1)
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
	if b, _ := s.CheckActorByName(a.Name); !b {
		pc, errP := s.db.Prepare(`INSERT INTO actors (name, sex, birthday) VALUES ($1, $2, $3)`)
		if errP != nil {
			return errP
		}
		_, err := pc.Exec(a.Name, a.Sex, a.Birthday)
		if err != nil {
			return err
		}
		return nil
	} else {
		return ErrAlreadyIn
	}
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
							SET name = $1,
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

func (s Storage) GetActor(id int) (Actor, error) {
	var act Actor
	row := s.db.QueryRow("SELECT name, sex, birthday FROM actors WHERE actor_id=$1", id)
	err := row.Scan(&act.Name, &act.Sex, &act.Birthday)
	if err != nil {
		return Actor{}, err
	}
	return act, nil
}

func (s Storage) CheckActorByName(name string) (bool, error) {
	rows, err := s.db.Query("SELECT actor_id FROM actors WHERE name = $1", name)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	// Обрабатываем каждую запись
	if rows.Next() {
		return true, nil
	} else {
		return false, nil
	}
}

func (s Storage) GetActors() ([]Actor, error) {
	var actors = []Actor{}
	rows, errQ := s.db.Query("SELECT actor_id, name, sex, birthday FROM actors")
	if errQ != nil {
		return []Actor{}, errQ
	}
	for rows.Next() {
		actor := Actor{}
		rows.Scan(&actor.ID, &actor.Name, &actor.Sex, &actor.Birthday)
		films, err := s.GetFilmsBySearch(actor.Name)
		if err == nil {
			actor.Films = films
		}
		actors = append(actors, actor)
	}
	return actors, nil
}

func (s Storage) SetFilm(f Film) error {
	if b, _ := s.CheckFilmByTitle(f.Title); !b {
		pc, errP := s.db.Prepare(`INSERT INTO films (title, description, date, rating, actors_list) VALUES ($1, $2, $3, $4, $5)`)
		if errP != nil {

			return errP
		}
		_, err := pc.Exec(f.Title, f.Description, f.Date, f.Rating, f.Actors)
		if err != nil {

			return err
		}
		return nil
	} else {
		return ErrAlreadyIn
	}

}

func (s Storage) DeleteFilm(id int) (int, error) {
	pc, errP := s.db.Prepare(`DELETE FROM films WHERE film_id = $1`)
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

func (s Storage) UpdateFilm(film Film) (int, error) {
	pc, errP := s.db.Prepare(`UPDATE films
							SET title = $1,
							description = $2,
							date = $3, 
							rating = $4,
							actors_list = $5
							WHERE film_id = $6`)
	if errP != nil {
		return 0, errP
	}
	r, err := pc.Exec(film.Title, film.Description, film.Date, film.Rating, film.Actors, film.ID)
	if err != nil {
		return 0, err
	}
	count, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (s Storage) GetFilmsBySort(field string, method string) ([]Film, error) {
	films := []Film{}
	pr :=  fmt.Sprintf(`SELECT film_id, title, description, date, rating, actors_list FROM films ORDER BY %s %s`, field, method)
	rows, err := s.db.Query(pr)
	if err != nil {
		return  []Film{}, err
	}

	for rows.Next() {
		film := Film{}
		err = rows.Scan(&film.ID, &film.Title, &film.Description, &film.Date, &film.Rating, &film.Actors)
		if err == nil {
			films = append(films,  film)
		}	
	}
	return films, nil
}

func (s Storage) CheckFilmByTitle(title string) (bool, error) {
	rows, err := s.db.Query("SELECT film_id FROM films WHERE title = $1", title)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	// Обрабатываем каждую запись
	if rows.Next() {
		return true, nil
	} else {
		return false, nil
	}
}
