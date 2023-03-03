package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time2do/database"
	"time2do/entity"

	_ "time2do/docs"

	"github.com/gorilla/mux"
	"github.com/thanhpk/randstr"
)

// @Summary 유저 생성하기
// @Tag User
// @Accept  json
// @Produce  json
// @Router /user [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := io.ReadAll(r.Body)
	var command CreateUserCommand
	_ = json.Unmarshal(requestBody, &command)

	var user *entity.User
	database.Connector.
		Select("id").
		Where(&entity.User{UserName: command.UserName}).
		Find(&user)

	if user.Id != nil {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(fmt.Sprintf("이미 존재하는 유저 이름입니다. userName : %s", command.UserName))
		return
	}

	user = &entity.User{
		IdToken:    randstr.Hex(4),
		UserName:   command.UserName,
		Password:   command.Password,
		Onboarding: false,
	}
	database.Connector.Create(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(*user)
}

type CreateUserCommand struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// @Summary 유저 ID 로 조회하기
// @Tag User
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
// @Tag User
// @Accept  json
// @Produce  json
// @Router /users [get]
func GetAllUser(w http.ResponseWriter, _ *http.Request) {
	var users []entity.User
	database.Connector.Find(&users)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	uintUserId, _ := strconv.ParseUint(userId, 10, 32)
	id := uint(uintUserId)

	var command updateUserCommand
	requestBody, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(requestBody, &command)

	user := entity.User{Id: &id}
	database.Connector.First(&user)

	if command.UserName != nil {
		user.UserName = *command.UserName
		database.Connector.Model(user).Update("user_name", &command.UserName)
	}

	if command.OnBoarding != nil {
		user.Onboarding = *command.OnBoarding
		database.Connector.Model(user).Update("onboarding", &command.OnBoarding)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

type updateUserCommand struct {
	UserName   *string `json:"userName"`
	OnBoarding *bool   `json:"onBoarding"`
}

type UserCommand struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// @Summary User Login
// @Tag User
// @Accept json
// @Produce json
// @Param body body UserCommand true "User credentials"
// @Router /login [post]
func LoginUser(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := io.ReadAll(r.Body)
	fmt.Printf("Request body: %s\n", string(requestBody))

	var command UserCommand
	_ = json.Unmarshal(requestBody, &command)

	var user *entity.User
	database.Connector.Where(&entity.User{UserName: command.UserName}).First(&user)

	if user.Id == nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Message: "존재하지 않는 유저입니다."})
		return
	}

	if user.Password != command.Password {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode("잘못된 비밀번호 입니다.")
		return
	}

	// TODO: Create and return an authentication token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}
