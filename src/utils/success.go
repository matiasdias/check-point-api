package utils

import (
	"encoding/json"
	"net/http"
)

// APISuccess representa um erro HTTP personalizado.
type APISuccess struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

// NewAPISuccess cria um novo erro HTTP com a mensagem e o código de status especificados.
func NewAPISuccess(message string, statusCode int) *APISuccess {
	return &APISuccess{
		Message:    message,
		StatusCode: statusCode,
	}
}

// RespondWithSuccess responde com um código de status HTTP 200 OK.
func RespondWithSuccess(w http.ResponseWriter, success *APISuccess) {
	responseJSON, marsharSuccess := json.Marshal(success)
	if marsharSuccess != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func RespondWithSuccessJSON(w http.ResponseWriter, data interface{}) {
	responseJSON, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Erro formatando resposta JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func RespondWithSuccessMessage(w http.ResponseWriter, message string) {
	success := NewAPISuccess(message, http.StatusOK)
	RespondWithSuccess(w, success)
}
