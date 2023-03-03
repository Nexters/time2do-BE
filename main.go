package main

import (
	"log"
	"net/http"
	"os"
	"time2do/controller"
	"time2do/database"

	_ "github.com/go-sql-driver/mysql" //Required for MySQL dialect
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

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

	router.HandleFunc("/users", controller.GetAllUser).Methods("GET")  // swg 0
	router.HandleFunc("/users", controller.CreateUser).Methods("POST") // swg 0
	router.HandleFunc("/users/{id}", controller.GetUserByID).Methods("GET")
	router.HandleFunc("/users/{id}", controller.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}/timeRecords", controller.SyncTimeRecords).Methods("POST")
	router.HandleFunc("/users/{id}/reports", controller.ViewReport).Methods("GET")

	router.HandleFunc("/login", controller.LoginUser).Methods("POST")

	router.HandleFunc("/timers", controller.CreateTimer).Methods("POST")
	router.HandleFunc("/timers", controller.GetAllTimers).Methods("GET")
	router.HandleFunc("/timers/{invitationCode}", controller.GetGroupTimer).Methods("GET")
	router.HandleFunc("/timers/{invitationCode}/supporting", controller.GetSupporting).Methods("GET")
	router.HandleFunc("/timers/{invitationCode}/users", controller.GetCountdownParticipants).Methods("GET")
	router.HandleFunc("/timers/{invitationCode}/users/{userId}", controller.Participate).Methods("POST")
	router.HandleFunc("/timers/{invitationCode}/users/{userId}", controller.Leave).Methods("DELETE")
	router.HandleFunc("/timers/{invitationCode}/users/{userId}/supporting", controller.MakeSupporting).Methods("POST")
	router.HandleFunc("/users/{userId}/timers/{timerId}/timeRecords", controller.CreateTimerRecord).Methods("POST")

	router.HandleFunc("/users/{userId}/tasks", controller.CreateToDos).Methods("POST")
	router.HandleFunc("/users/{userId}/tasks", controller.GetToDoById).Methods("GET")
	router.HandleFunc("/tasks", controller.GetAllToDo).Methods("GET")
	router.HandleFunc("/tasks/{id}", controller.GetToDoById).Methods("GET")

	// Swagger
	corsHandler := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}),
	)
	router.PathPrefix("/swagger").Handler(corsHandler(httpSwagger.WrapHandler))
}

// @title Swagger Time2Do API
// @version 1.0
// @description This is Time2Do API Server
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email devgunho@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8888
// @BasePath /
func main() {
	initDB()

	log.Println("[*] Starting the HTTP server on port 8888")
	router := mux.NewRouter().StrictSlash(true)
	initHandlers(router)

	corsHandler := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}),
	)
	log.Fatal(http.ListenAndServe(":8888", corsHandler(router)))
}
