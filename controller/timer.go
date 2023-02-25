package controller

import (
	"crypto/rand"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"time2do/database"
	"time2do/entity"
)

const otpChars = "1234567890"

// @Summary 타이머 생성하기
// @Accept  json
// @Produce  json
// @Router /group [post]
func CreateTimer(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := io.ReadAll(r.Body)
	var timer entity.Timer
	_ = json.Unmarshal(requestBody, &timer)
	if timer.Type == entity.Group {
		otp, _ := generateOTP(6)
		timer.InvitationCode = &otp
	}

	database.Connector.Create(timer)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(timer)
}

// @Summary 그룹 조회하기
// @Accept  json
// @Produce  json
// @Router /groups [get]
func GetAllTimers(w http.ResponseWriter, _ *http.Request) {
	var timers []entity.Timer
	database.Connector.Find(&timers)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(timers)
}

func CreateTimerRecord(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := io.ReadAll(r.Body)

	vars := mux.Vars(r)
	userId := vars["userId"]
	timerId := vars["timerId"]

	uIntUserId, _ := strconv.ParseUint(userId, 10, 32)
	uIntTimerId, _ := strconv.ParseUint(timerId, 10, 32)

	var createTimeRecordCommand CreateTimeRecordCommand
	_ = json.Unmarshal(requestBody, &createTimeRecordCommand)

	// TODO: 동일한 시간에 다른 타이머 기록이 있을 경우 예외 발생
	timeRecord := entity.TimeRecord{TimerId: uint(uIntTimerId), UserId: uint(uIntUserId), StartTime: createTimeRecordCommand.StartTime, EndTime: createTimeRecordCommand.EndTime}
	database.Connector.Create(timeRecord)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(timeRecord)
}

func GetCountdownParticipants(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	timerId := vars["timerId"]
	uIntTimerId, _ := strconv.ParseUint(timerId, 10, 32)

	var timer entity.Timer

	database.Connector.
		Where(&entity.Timer{Id: uint(uIntTimerId)}).
		Preload("Users").
		Find(&timer)

	var participants []Participant
	for _, participant := range timer.Users {
		participants = append(participants, Participant{UserName: participant.UserName})
	}

	_ = json.NewEncoder(w).Encode(participants)
}

func Participate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	timerId := vars["timerId"]

	uIntUserId, _ := strconv.ParseUint(userId, 10, 32)
	uIntTimerId, _ := strconv.ParseUint(timerId, 10, 32)

	var timer entity.Timer

	database.Connector.
		Where(&entity.Timer{Id: uint(uIntTimerId)}).
		Preload("Users").
		Find(&timer)

	for _, user := range timer.Users {
		if uint(uIntUserId) == user.Id {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(409)
			return
		}
	}

	user := entity.User{Id: uint(uIntUserId)}
	database.Connector.First(&user)
	timer.Users = append(timer.Users, &user)
	database.Connector.Updates(timer)
}

type Participant struct {
	UserName string `json:"userName"`
}

type CreateTimeRecordCommand struct {
	StartTime entity.DateTime `json:"startDateTime"`
	EndTime   entity.DateTime `json:"endDateTime"`
}

func generateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}
