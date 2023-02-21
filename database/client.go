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
<<<<<<< HEAD

// Migrate create/updates database table
func UserMigrate(table *entity.User) {
	Connector.AutoMigrate(&table)
	log.Println("[+] MySQL: User Table migrated")
}

func GroupMigrate(table *entity.Timer) {
	Connector.AutoMigrate(&table)
	log.Println("[+] MySQL: Group Table migrated")
}

func TaskMigrate(table *entity.ToDo) {
	Connector.AutoMigrate(&table)
	log.Println("[+] MySQL: Task Table migrated")
}

func ParticipateMigrate(table *entity.Participant) {
	Connector.AutoMigrate(&table)
	log.Println("[+] MySQL: Participate Table migrated")
}
=======
>>>>>>> origin/feature
