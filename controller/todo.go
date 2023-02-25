package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
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
	requestBody, _ := io.ReadAll(r.Body)
	var command CreateToDoCommand
	_ = json.Unmarshal(requestBody, &command)

	vars := mux.Vars(r)
	userId := vars["userId"]
	uIntUserId, _ := strconv.ParseUint(userId, 10, 32)
	id := uint(uIntUserId)

	toDo := entity.ToDo{
		UserId:        id,
		Content:       command.Content,
		Completed:     command.Completed,
		CreatedTime:   command.CreatedTime,
		CompletedTime: command.CompletedTime,
	}
	database.Connector.Create(&toDo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(toDo)
}

type CreateToDoCommand struct {
	Content       string    `json:"content"`
	Completed     bool      `json:"completed"`
	CreatedTime   DateTime  `json:"createdTime"`
	CompletedTime *DateTime `json:"completedTime"`
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
