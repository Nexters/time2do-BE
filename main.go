package main

import (
	"log"
	"net/http"
	"os"
	"time2do/controller"
	"time2do/database"
	"time2do/entity"

	_ "github.com/go-sql-driver/mysql" //Required for MySQL dialect
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func initDB() {
	log.Println("[*] MySQL DB Setting...")

	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	config := database.Config{
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		DB:       os.Getenv("MYSQL_DATABASE"),
	}

	connectionString := database.GetConnectionString(config)
	log.Println("[+] Connection String Check:", connectionString)
	connectErr := database.Connect(connectionString)
	if connectErr != nil {
		panic(connectErr.Error())
	}
	database.Migrate(&entity.User{})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	log.Println(r.Header)
	w.Write([]byte("Time2Do Server is healthy."))
}

func initHandlers(router *mux.Router) {
	router.HandleFunc("/", healthCheck).Methods("GET")
	router.HandleFunc("/user", controller.CreateUser).Methods("POST")
	router.HandleFunc("/users", controller.GetAllUser).Methods("GET")

}

func main() {
	initDB()

	log.Println("[*] Starting the HTTP server on port 8888")
	router := mux.NewRouter().StrictSlash(true)
	initHandlers(router)
	log.Fatal(http.ListenAndServe(":8888", router))
}
