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

func CreateToDos(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := io.ReadAll(r.Body)
	var commands []CreateToDoCommand
	_ = json.Unmarshal(requestBody, &commands)

	if len(commands) == 0 {
		return
	}

	vars := mux.Vars(r)
	userId := vars["userId"]
	uIntUserId, _ := strconv.ParseUint(userId, 10, 32)
	id := uint(uIntUserId)

	var toDos []entity.ToDo
	for _, command := range commands {
		toDos = append(
			toDos,
			entity.ToDo{
				UserId:        id,
				Content:       command.Content,
				Completed:     command.Completed,
				CreatedTime:   command.CreatedTime,
				CompletedTime: command.CompletedTime,
			},
		)
	}

	database.Connector.CreateInBatches(&toDos, len(toDos))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(toDos)
}

type CreateToDoCommand struct {
	Content       string    `json:"content"`
	Completed     bool      `json:"completed"`
	CreatedTime   DateTime  `json:"createdTime"`
	CompletedTime *DateTime `json:"completedTime"`
}

func GetAllToDo(w http.ResponseWriter, r *http.Request) {
	var tasks []entity.ToDo
	database.Connector.Find(&tasks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(tasks)
}

// func GetToDoById(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	key := vars["id"]
// 	var task entity.ToDo
// 	database.Connector.First(&task, key)
// 	w.Header().Set("Content-Type", "application/json")
// 	_ = json.NewEncoder(w).Encode(task)
// }

func GetToDosByUserId(userId uint) ([]entity.ToDo, error) {
	var toDos []entity.ToDo
	if err := database.Connector.Where("user_id = ?", userId).Find(&toDos).Error; err != nil {
		return nil, err
	}
	return toDos, nil
}

func GetToDoById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.ParseUint(vars["userId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	toDos, err := GetToDosByUserId(uint(userId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(toDos)
}
