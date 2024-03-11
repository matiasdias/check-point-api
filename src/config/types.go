package config

// DatabaseConfig estrutura do banco de dados
type DatabaseConfig struct {
	DBName     string `json:"db_name"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
}

// APIConfig estrutura da API
type APIConfig struct {
	APIPort int `json:"api_port"`
}
