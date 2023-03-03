package database

import (
	"fmt"
	"net/url"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	DB       string
}

var GetConnectionString = func(config Config) string {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s", config.User, config.Password, config.Host, config.Port, config.DB, url.PathEscape("Asia/Seoul"))
	return connectionString
}
