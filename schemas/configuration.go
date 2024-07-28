package schemas

type DBConfig struct {
	DBName     string `json:"db_name"`
	DBPassword string `json:"db_password"`
	DBUser     string `json:"db_user"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DBSSLMode  string `json:"db_ssl_mode"`
}

type ProjectConfiguration struct {
	Environment string   `json:"environment"`
	DBConfig    DBConfig `json:"database"`
	DBString    string   `json:"database_string,omitempty"`
}