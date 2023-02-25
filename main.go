package main

import (
	_ "github.com/go-sql-driver/mysql" //Required for MySQL dialect
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time2do/controller"
	"time2do/database"

	httpSwagger "github.com/swaggo/http-swagger"
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
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	log.Println(r.Header)
	_, _ = w.Write([]byte("Time2Do Server is healthy."))
}

func initHandlers(router *mux.Router) {
	router.HandleFunc("/", healthCheck).Methods("GET")

	router.HandleFunc("/users", controller.GetAllUser).Methods("GET")
	router.HandleFunc("/users", controller.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", controller.GetUserByID).Methods("GET")
	router.HandleFunc("/users/{id}", controller.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}/reports", controller.ViewReport).Methods("GET")

	router.HandleFunc("/timers", controller.CreateTimer).Methods("POST")
	router.HandleFunc("/timers", controller.GetAllTimers).Methods("GET")
	router.HandleFunc("/timers/{timerId}/users", controller.GetCountdownParticipants).Methods("GET")
	router.HandleFunc("/timers/{timerId}/users/{userId}", controller.Participate).Methods("POST")
	router.HandleFunc("/users/{userId}/timers/{timerId}/timeRecords", controller.CreateTimerRecord).Methods("POST")

	router.HandleFunc("/tasks", controller.CreateTask).Methods("POST")
	router.HandleFunc("/tasks", controller.GetAllTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", controller.GetTaskByID).Methods("GET")

	// Swagger
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
}

// @title Orders API
// @version 1.0
// @description This is a sample serice for managing orders
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email soberkoder@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8888
// @BasePath /
func main() {
	initDB()

	log.Println("[*] Starting the HTTP server on port 8888")
	router := mux.NewRouter().StrictSlash(true)
	initHandlers(router)
	log.Fatal(http.ListenAndServe(":8888", handlers.CORS()(router)))
}
