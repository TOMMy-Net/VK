package services

import (
	"fmt"
	"os"
	"time"

	"github.com/TOMMy-Net/VK/db"

	"github.com/golang-jwt/jwt/v5"
)

type AuthStore interface {
	GetUser(name, password string) (db.User, error)
	GetTokenByID(id int) (string, error)
}

type AuthService struct {
	store AuthStore
}

func NewAuthService(s AuthStore) *AuthService {
	return &AuthService{store: s}
}

func (a *AuthService) SignInUser(username, password string) (string, error) {
	user, err := a.store.GetUser(username, password)
	if err != nil {
		return "", err
	}

	// Создание claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": user.Role,
		"id":   user.ID,
		"exp":  time.Now().Add(time.Hour * 720).Unix(),
	})

	// Подписание токена
	token, err := claims.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return "", err
	}
	return token, nil
}

func (a *AuthService) AuthUser(tokenString string) (jwt.MapClaims, error) {
	// Парсинг токена
	t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	// Проверка валидности токена
	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {

		if int(time.Unix(int64(claims["exp"].(float64)), 0).Unix()) > int(time.Now().Unix()) {

			return claims, nil

		}
		return nil, nil
	} else {
		fmt.Println("erer")
		return nil, nil
	}

}
