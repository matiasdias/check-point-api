package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	Door      = 0
	Driver    = "postgres"
	SecretKey []byte
)

func Connection() (*sql.DB, error) {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	Door, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		Door = 3001
	}

	dbUser := os.Getenv("DB_USER")
	pgPass := os.Getenv("DB_PASSWORD")
	pgHost := os.Getenv("DB_HOST")
	pgPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	SecretKey = []byte(os.Getenv("API_SECRET"))

	Connection := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		pgHost, pgPort, dbUser, dbName, pgPass,
	)

	db, err := sql.Open(Driver, Connection)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		defer db.Close()
		return nil, err
	}

	return db, nil

}
