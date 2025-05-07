package config

import (
	"github.com/go-playground/validator/v10"
	"log"
	"strings"
)

func NewValidator() *validator.Validate {
	var err error
	validation := validator.New()

	// custom validation here
	err = validation.RegisterValidation("not_only_space", NotOnlySpace)
	if err != nil {
		log.Println("[ERROR] Error register validation not_only_space : " + err.Error())
	}

	return validation
}

func NotOnlySpace(fl validator.FieldLevel) bool {
	data := fl.Field().String()

	return strings.TrimSpace(data) != ""
}
