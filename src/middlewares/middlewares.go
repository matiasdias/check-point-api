package middlewares

import (
	"check-point/src/token"
	"check-point/src/utils"
	"log"
	"net/http"
	"time"
)

// Logger
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now()
		formattedTime := currentTime.Format("2006-01-02 15:04:05.999999")

		log.Printf("\n%s %q %q %q", formattedTime, r.Host, r.RequestURI, r.Method)
		next(w, r)
	}
}

// Authenticate verifica se o usuario está autenticado e valida o token
func Authenticate(ProximaFuncao http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if erro := token.ValidateToken(r); erro != nil {
			err := utils.UnauthorizedError("Token contains as invalid number")
			utils.RespondWithError(w, err)
			return
		}
		ProximaFuncao(w, r)
	}
}
