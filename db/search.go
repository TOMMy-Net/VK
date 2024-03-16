package db

import (
	"strings"

	_ "github.com/lib/pq"
)

func (s Storage) GetFilmsBySearch(name string) ([]Film, error) {
	var films = []Film{}
	rows, errQ := s.db.Query("SELECT film_id, title, description, date, rating, actors_list FROM films ")
	if errQ != nil {
		return []Film{}, errQ
	}
	for rows.Next() {
		film := Film{}
		rows.Scan(&film.ID, &film.Title, &film.Description, &film.Date, &film.Rating, &film.Actors)
	
		if strings.Contains(strings.ToLower(film.Title), strings.ToLower(name)) || strings.Contains(strings.ToLower(film.Actors), strings.ToLower(name)) {
			films = append(films, film)
		}
	}
	return films, nil
}

