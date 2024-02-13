package config

import (
	"database/sql"
	"encoding/json"
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

type DatabaseConfig struct {
	DBName     string `json:"db_name"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
}

func LoadDatabaseConfig(filename string) (DatabaseConfig, error) {
	var dbConfig DatabaseConfig
	file, err := os.Open(filename)
	if err != nil {
		return dbConfig, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&dbConfig)
	if err != nil {
		return dbConfig, fmt.Errorf("failed to decode config file: %w", err)
	}
	return dbConfig, nil
}

func Connection() (*sql.DB, error) {
	var err error
	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	apiPortStr := os.Getenv("API_PORT")
	Door, err = strconv.Atoi(apiPortStr)
	if err != nil {
		Door = 3002
	}

	dbConfig, err := LoadDatabaseConfig("config/config.api.json")
	if err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	SecretKey = []byte(os.Getenv("API_SECRET"))

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBUser, dbConfig.DBName, dbConfig.DBPassword,
	)

	db, err := sql.Open(Driver, connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	//log.Printf("Successfully connected to database at %s:%s", dbConfig.DBHost, dbConfig.DBPort)

	return db, nil

}
