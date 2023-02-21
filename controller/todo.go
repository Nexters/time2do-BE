package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"time2do/database"
	"time2do/entity"

	"github.com/gorilla/mux"
)

// @Summary 할일 생성하기
// @Tags ToDo (Task)
// @Accept  json
// @Produce  json
// @Router /task [post]
<<<<<<< HEAD:controller/task.go
func CreateToDo(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var task entity.ToDo
	json.Unmarshal(requestBody, &task)
=======
func CreateTask(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := io.ReadAll(r.Body)
	var task entity.ToDo
	_ = json.Unmarshal(requestBody, &task)
>>>>>>> origin/feature:controller/todo.go
	database.Connector.Create(task)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(task)
}

// @Summary 아무 조건 없이 모든 ToDo 불러오기
// @Tags ToDo (Task)
// @Accept  json
// @Produce  json
// @Router /tasks [get]
<<<<<<< HEAD:controller/task.go
func GetAllToDo(w http.ResponseWriter, r *http.Request) {
=======
func GetAllTask(w http.ResponseWriter, r *http.Request) {
>>>>>>> origin/feature:controller/todo.go
	var tasks []entity.ToDo
	database.Connector.Find(&tasks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(tasks)
}

<<<<<<< HEAD:controller/task.go
// @Summary ToDo Id 를 통해 ToDo 불러오기
=======
// @Summary ToDo ID 를 통해 ToDo 불러오기
>>>>>>> origin/feature:controller/todo.go
// @Tags ToDo (Task)
func GetToDoById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	var task entity.ToDo
	database.Connector.First(&task, key)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(task)
}
