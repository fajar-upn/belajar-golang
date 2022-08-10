package helper

import "github.com/go-playground/validator/v10"

/**
this helper for add 'meta' respone in database
so, user can receive success, error, or failed
*/

type Response struct {
	Meta Meta
	Data interface{} //interface{}, because data will return flexible value
}

type Meta struct {
	Message string
	Code    int
	Status  string
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}

func FormatValidationError(err error) []string {
	var errors []string //this code for accomodate error from 'err' json

	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}

	return errors
}
