package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var Connector *gorm.DB

// Connect creates MySQL connection
func Connect(connectionString string) error {
	var err error
	Connector, err = gorm.Open(mysql.Open(connectionString))
	if err != nil {
		return err
	}
	log.Println("[+] Connection was successful")
	return nil
}
