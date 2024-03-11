package utils

import (
	"encoding/json"
	"net/http"
)

// APIError representa um erro HTTP personalizado.
type APIError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

// NewAPIError cria um novo erro HTTP com a mensagem e o código de status especificados.
func NewAPIError(message string, statusCode int) *APIError {
	return &APIError{
		Message:    message,
		StatusCode: statusCode,
	}
}

// RespondWithError responde com um erro HTTP contendo a mensagem formatada em JSON.
func RespondWithError(w http.ResponseWriter, err *APIError) {
	responseJSON, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)
	w.Write(responseJSON)
}

// UnauthorizedError retorna um erro HTTP 401 Unauthorized.
func UnauthorizedError(message string) *APIError {
	return NewAPIError(message, http.StatusUnauthorized)
}

// BadRequestError retorna um erro HTTP 400 Bad Request.
func BadRequestError(message string) *APIError {
	return NewAPIError(message, http.StatusBadRequest)
}

// NotFoundError retorna um erro HTTP 404 Not Found.
func NotFoundError(message string) *APIError {
	return NewAPIError(message, http.StatusNotFound)
}

// InternalServerError retorna um erro HTTP 500 Internal Server Error.
func InternalServerError(message string) *APIError {
	return NewAPIError(message, http.StatusInternalServerError)
}

// ForbiddenError retorna um erro HTTP 403 Forbidden.
func ForbiddenError(message string) *APIError {
	return NewAPIError(message, http.StatusForbidden)
}

// ConflitError retorna um erro HTTP 409 Conflit.
func ConflitError(message string) *APIError {
	return NewAPIError(message, http.StatusConflict)
}

// UnprocessableEntityError retorna um erro HTTP 422 Unprocessable Entity.
func UnprocessableEntityError(message string) *APIError {
	return NewAPIError(message, http.StatusUnprocessableEntity)
}
