package db

import (
	"github.com/go-playground/validator/v10"
)

func ValidStruct(f interface{}) error {
	validate := validator.New()
	err := validate.Struct(f)
	if err != nil{
		return err
	}
	return nil
}