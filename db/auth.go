package db

import (
	"github.com/TOMMy-Net/VK/internal"
	_ "github.com/lib/pq"
)

func (s *Storage) GetUser(username, password string) (User, error) {
	var user User
	row := s.db.QueryRow("SELECT user_id, name, password, role FROM users WHERE name=$1 AND password=$2", username, internal.Hash([]byte(password)))
	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.Role)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *Storage) CheckUser(username string) (bool, error) {
	rows, err := s.db.Query(`SELECT user_id FROM users WHERE name = $1`, username)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if rows.Next() {
		return true, nil
	} else {
		return false, nil
	}
}

func (s *Storage) CreateUser(u User) error {
	if b, _ := s.CheckUser(u.Name); !b {
		pc, errP := s.db.Prepare(`INSERT INTO users (name, password) VALUES ($1, $2)`)
		if errP != nil {
			return errP
		}
		_, err := pc.Exec(u.Name, internal.Hash([]byte(u.Password)))
		if err != nil {
			return err
		}
		return nil
	} else {
		return ErrAlreadyIn
	}
}

func (s *Storage) UpdateToken(token string, id int) (int, error) {
	pc, errP := s.db.Prepare(`UPDATE users SET token = $1 WHERE user_id = $2`)
	if errP != nil {
		return 0, errP
	}
	r, err := pc.Exec(token, id)

	if err != nil {
		return 0, err
	}
	count, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (s *Storage) GetTokenByID(id int) (string, error) {
	var token string
	row := s.db.QueryRow("SELECT token FROM users WHERE user_id=$1", id)
	err := row.Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}
