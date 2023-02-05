package database

import (
	"log"
	"time2do/entity"

	"github.com/jinzhu/gorm"
)

var Connector *gorm.DB

// Connect creates MySQL connection
func Connect(connectionString string) error {
	var err error
	Connector, err = gorm.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	log.Println("[+] Connection was successful")
	return nil
}

// Migrate create/updates database table
func UserMigrate(table *entity.User) {
	Connector.AutoMigrate(&table)
	log.Println("[+] MySQL: User Table migrated")
}

func GroupMigrate(table *entity.Group) {
	Connector.AutoMigrate(&table)
	log.Println("[+] MySQL: Group Table migrated")
}
