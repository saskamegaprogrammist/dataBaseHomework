package models

type RequestError struct {
	Message string `json:"message"`
}

func CreateError(message string) RequestError {
	return RequestError{Message:message}
}