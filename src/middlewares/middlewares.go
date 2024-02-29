package middlewares

import (
	"check-point/src/token"
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

// Autenticar verifica se o usuario fazendo a requisição está autenticado
func Authenticate(ProximaFuncao http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if erro := token.ValidateToken(r); erro != nil {
			http.Error(w, "token contains as invalid number", http.StatusUnauthorized)
			return
		}
		ProximaFuncao(w, r)
	}
}
