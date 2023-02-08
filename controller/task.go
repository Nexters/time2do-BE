package controller

import (
	"encoding/json"
	"io/ioutil"
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
func CreateTask(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var task entity.Task
	json.Unmarshal(requestBody, &task)
	database.Connector.Create(task)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// @Summary 아무 조건 없이 모든 Task 불러오기
// @Tags ToDo (Task)
// @Accept  json
// @Produce  json
// @Router /tasks [get]
func GetAllTask(w http.ResponseWriter, r *http.Request) {
	var tasks []entity.Task
	database.Connector.Find(&tasks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

// @Summary Task ID 를 통해 Task 불러오기
// @Tags ToDo (Task)
func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	var task entity.Task
	database.Connector.First(&task, key)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}
