package middlewares

import (
	"check-point/src/auth"
	"log"
	"net/http"
	"time"
)

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now()

		formattedTime := currentTime.Format("2006-01-02 15:04:05.999999")

		log.Printf("\n%s %q %q %q", formattedTime, r.Host, r.RequestURI, r.Method)

		//body, err := io.ReadAll(r.Body)
		//if err != nil {
		//	log.Printf("Erro ao ler o corpo da requisição: %v", err)
		//} else {
		//	log.Printf("Corpo da requisição: %s", body)
		//}
		//r.Body = io.NopCloser(bytes.NewBuffer(body))

		next(w, r)
	}
}

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := auth.ValidateToken(r); err != nil {
			http.Error(w, "token contains as invalid number", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
