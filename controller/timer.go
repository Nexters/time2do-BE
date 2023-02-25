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
	var command createTimerCommand
	_ = json.Unmarshal(requestBody, &command)
	var timer entity.Timer
	timer = entity.Timer{
		Name:      command.Name,
		MakerId:   command.MakerId,
		Type:      command.Type,
		Tags:      command.Tag,
		StartTime: command.StartTime,
		EndTime:   command.EndTime,
	}
	if command.Type == entity.Group {
		otp, _ := generateOTP(6)
		timer.InvitationCode = &otp
	}

	database.Connector.Create(&timer)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(timer)
}

type createTimerCommand struct {
	Name      string           `json:"name"`
	MakerId   uint             `json:"makerId"`
	Type      entity.TimerType `json:"type"`
	Tag       string           `json:"tag"`
	StartTime DateTime         `json:"startTime"`
	EndTime   *DateTime        `json:"endTime"`
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

func GetGroupTimer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invitationCode := vars["invitationCode"]

	var timer entity.Timer
	database.Connector.Where(entity.Timer{InvitationCode: &invitationCode}).First(&timer)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(timer)
}

func GetCountdownParticipants(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invitationCode := vars["invitationCode"]

	var timer entity.Timer
	database.Connector.
		Where(&entity.Timer{InvitationCode: &invitationCode}).
		Preload("Users").
		Find(&timer)

	var participants []Participant
	for _, participant := range timer.Users {
		var toDos []entity.ToDo
		// TODO: private
		database.Connector.Where(&entity.ToDo{UserId: *participant.Id}).
			Find(&toDos)
		participants = append(participants, Participant{UserName: participant.UserName, ToDos: toDos})
	}

	_ = json.NewEncoder(w).Encode(participants)
}

func Participate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	invitationCode := vars["invitationCode"]

	uIntUserId, _ := strconv.ParseUint(userId, 10, 32)

	var timer entity.Timer

	database.Connector.
		Where(&entity.Timer{InvitationCode: &invitationCode}).
		Preload("Users").
		Find(&timer)

	id := uint(uIntUserId)
	contains := false
	for _, user := range timer.Users {
		if id == *user.Id {
			contains = true
			break
		}
	}
	if !contains {
		user := entity.User{Id: &id}
		database.Connector.First(&user)
		timer.Users = append(timer.Users, &user)
		database.Connector.Updates(timer)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(timer)
}

type Participant struct {
	UserName string        `json:"userName"`
	ToDos    []entity.ToDo `json:"toDos"`
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
