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
// @Tag ToDo (Task)
// @Accept  json
// @Produce  json
// @Router /task [post]
func CreateToDo(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var task entity.ToDo
	json.Unmarshal(requestBody, &task)
	database.Connector.Create(task)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(task)
}

// @Summary 아무 조건 없이 모든 ToDo 불러오기
// @Tag ToDo (Task)
// @Accept  json
// @Produce  json
// @Router /tasks [get]
func GetAllToDo(w http.ResponseWriter, r *http.Request) {
	var tasks []entity.ToDo
	database.Connector.Find(&tasks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(tasks)
}

// @Summary ToDo ID 를 통해 ToDo 불러오기
// @Tag ToDo (Task)
func GetToDoById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	var task entity.ToDo
	database.Connector.First(&task, key)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(task)
}
