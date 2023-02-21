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
	router.HandleFunc("/user", controller.CreateUser).Methods("POST")
	router.HandleFunc("/user/{id}", controller.GetUserByID).Methods("GET")
	router.HandleFunc("/users", controller.GetAllUser).Methods("GET")

	router.HandleFunc("/group", controller.CreateTimer).Methods("POST")
	router.HandleFunc("/groups", controller.GetAllTimers).Methods("GET")

	router.HandleFunc("/task", controller.CreateTask).Methods("POST")
	router.HandleFunc("/tasks", controller.GetAllTask).Methods("GET")
	router.HandleFunc("/task/{id}", controller.GetTaskByID).Methods("GET")

	router.HandleFunc("/report/{id}", controller.ViewReport).Methods("GET")

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
