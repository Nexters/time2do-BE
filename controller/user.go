package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time2do/database"
	"time2do/entity"

	_ "time2do/docs"

	"github.com/gorilla/mux"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var user entity.User
	json.Unmarshal(requestBody, &user)
	database.Connector.Create(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	var user entity.User
	database.Connector.First(&user, key)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// @Summary Get details of all users
// @Description Get details of all users
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.User
// @Router /users [get]
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	var users []entity.User
	database.Connector.Find(&users)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
