package db

import (
	"check-point/src/config"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	// Driver nome do driver
	Driver = "postgres"

	APIConfigInfo config.APIConfig

	// SecretKey chave de segurança
	SecretKey []byte
)

// loadAPIConfig carrega o arquivo de configuração da API
func loadAPIConfig(filePath string) (config.APIConfig, error) {
	var config config.APIConfig

	// Lê o conteúdo do arquivo JSON
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return config, err
	}

	// Decodifica o conteúdo do arquivo JSON para a estrutura APIConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return config, err
	}

	return config, nil
}

// loadDatabaseConfig carrega o arquivo de configuração do banco de dados
func loadDatabaseConfig(filePath string) (config.DatabaseConfig, error) {
	var dbConfig config.DatabaseConfig
	file, err := os.Open(filePath)
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

// Connection conecta com o banco de dados
func Connection() (*sql.DB, error) {

	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	SecretKey = []byte(os.Getenv("SECRET_KEY"))

	dbConfig, err := loadDatabaseConfig("config/config.api.json")
	if err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	APIConfigInfo, err = loadAPIConfig("config/config.api.json")
	if err != nil {
		return nil, fmt.Errorf("failed to load API config: %w", err)
	}

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

	return db, nil

}
