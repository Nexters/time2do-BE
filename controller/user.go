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

// @Summary 유저 생성하기
// @Tags user
// @Accept  json
// @Produce  json
// @Router /user [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var user entity.User
	json.Unmarshal(requestBody, &user)
	database.Connector.Create(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// @Summary 유저 ID 로 조회하기
// @Tags user
// @Accept  json
// @Produce  json
// @Router /user/{id} [get]
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	var user entity.User
	database.Connector.First(&user, key)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// @Summary 아무 조건 없이 모든 User 불러오기
// @Tags user
// @Accept  json
// @Produce  json
// @Router /users [get]
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	var users []entity.User
	database.Connector.Find(&users)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
