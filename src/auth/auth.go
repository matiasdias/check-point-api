package auth

import (
	config "check-point/src/db"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func CreateToken(employeeID uint64, employeeEmail string) (string, error) {
	permission := jwt.MapClaims{}
	permission["authorized"] = true
	permission["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permission["employeeID"] = employeeID
	permission["email"] = employeeEmail

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permission)

	secretKey := []byte(os.Getenv("API_SECRET"))
	return token.SignedString([]byte(secretKey))
}

func ValidateToken(r *http.Request) error {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, returnKeyVerify)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}
	return errors.New("Invalid token!")
}

func ExtractIDUserToken(r *http.Request) (uint64, error) {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, returnKeyVerify)
	if err != nil {
		return 0, err
	}
	if permission, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		employeeID, err := strconv.ParseUint(fmt.Sprintf("%.0f", permission["employeeID"]), 10, 64)
		if err != nil {
			return 0, err
		}
		return employeeID, nil
	}
	return 0, errors.New("Invalid token!")
}

func extractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}
	return ""
}

func returnKeyVerify(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Metodo de assinatura inesperado: %v", token.Header["alg"])
	}
	return config.SecretKey, nil
}
