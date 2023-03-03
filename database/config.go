package database

import (
	"fmt"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	DB       string
}

var GetConnectionString = func(config Config) string {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Asia/Seoul", config.User, config.Password, config.Host, config.Port, config.DB)
	return connectionString
}
