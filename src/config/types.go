package config

type DatabaseConfig struct {
	DBName     string `json:"db_name"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
}

type APIConfig struct {
	APIPort   int    `json:"api_port"`
	APISecret string `json:"api_secret"`
}
