package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time2do/database"
	"time2do/entity"

	_ "time2do/docs"

	"github.com/gorilla/mux"
	"github.com/thanhpk/randstr"

	"time"
)

// @Summary 유저 생성하기
// @Tags User
// @Accept  json
// @Produce  json
// @Router /user [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := io.ReadAll(r.Body)
	var reqUser entity.User
	_ = json.Unmarshal(requestBody, &reqUser)

	// logging for debug
	log.Println("\n" + string(requestBody))
	log.Println("User identify string: " + reqUser.IdToken)
	log.Println("Username: " + reqUser.UserName)
	log.Println("Password: " + reqUser.Password)

	w.Header().Set("Content-Type", "application/json")

	idToken := randstr.Hex(4)
	log.Println(idToken)

	var dbUser entity.User
	database.Connector.Where(&entity.User{IdToken: idToken}).Find(&dbUser)
	log.Println(dbUser.IdToken)

	reqUser.IdToken = idToken
	now := time.Now()
	reqUser.Id = uint(now.Nanosecond())
	log.Println("[2] User identify id: " + string(reqUser.Id))
	log.Println("[2] User identify string: " + reqUser.IdToken)
	log.Println("[2] Username: " + reqUser.UserName)
	log.Println("[2] Password: " + reqUser.Password)
	if results := database.Connector.Create(reqUser); results.Error != nil {
		w.WriteHeader(http.StatusConflict)
		// 	_ = json.NewEncoder(w).Encode("이미 존재하는  입니다")
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(reqUser)
}

// @Summary 유저 ID 로 조회하기
// @Tags User
// @Accept  json
// @Produce  json
// @Router /user/{id} [get]
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	var user entity.User
	database.Connector.First(&user, key)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}

// @Summary 아무 조건 없이 모든 User 불러오기
// @Tags User
// @Accept  json
// @Produce  json
// @Router /users [get]
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	var users []entity.User
	database.Connector.Find(&users)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}
