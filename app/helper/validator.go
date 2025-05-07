package helper

import "github.com/go-playground/validator/v10"

type ValidationErrorData struct {
	Field   string `json:"field"`
	Key     string `json:"key"`
	Tag     string `json:"tag"`
	Param   string `json:"param"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func ErrorValidationFormatter(errors validator.ValidationErrors) []ValidationErrorData {
	var formattedErrors []ValidationErrorData
	for _, fieldError := range errors {
		formattedErrors = append(formattedErrors, ValidationErrorData{
			Field:   fieldError.Field(),
			Key:     fieldError.Namespace(),
			Tag:     fieldError.Tag(),
			Param:   fieldError.Param(),
			Error:   fieldError.Error(),
			Message: ErrorValidationMessageGenerator(fieldError.Field(), fieldError.Tag(), fieldError.Param()),
		})
	}

	return formattedErrors
}

func ErrorValidationMessageGenerator(field, tag, param string) string {
	message := "Field validation for '" + field + "' failed on the '" + tag + "' tag"
	switch tag {
	case "required":
		message = field + " cannot be empty."
	case "not_only_space":
		message = field + " value cannot be only 'space'."
	}

	return message
}
