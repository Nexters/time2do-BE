package controller

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
	"time2do/database"
	"time2do/entity"

	"gorm.io/gorm"

	"github.com/gorilla/mux"
)

const otpChars = "1234567890"

var supportingMap = map[string][]supporting{}

func CreateTimer(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := io.ReadAll(r.Body)
	var command createTimerCommand
	_ = json.Unmarshal(requestBody, &command)
	var timer entity.Timer
	timer = entity.Timer{
		Name:      command.Name,
		MakerId:   command.MakerId,
		Type:      command.Type,
		Tag:       command.Tag,
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
	Tag       *string          `json:"tag"`
	StartTime DateTime         `json:"startTime"`
	EndTime   *DateTime        `json:"endTime"`
}

// @Summary 타이머 전부 조회
// @Tags Timer
// @Accept json
// @Produce json
// @Router /timers [get]
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

	//goland:noinspection GoPreferNilSlice
	var participants = []Participant{}
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

	if err := database.Connector.
		Where(&entity.Timer{InvitationCode: &invitationCode}).
		Preload("Users").
		Find(&timer).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

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

func Leave(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	invitationCode := vars["invitationCode"]
	uIntUserId, _ := strconv.ParseUint(userId, 10, 32)
	id := uint(uIntUserId)

	var timer entity.Timer
	database.Connector.
		Where(&entity.Timer{InvitationCode: &invitationCode}).
		Preload("Users").
		Find(&timer)

	users := timer.Users
	for _, user := range users {
		if id == *user.Id {
			_ = database.Connector.Model(&timer).Association("Users").Delete(user)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(timer)
			return
		}
	}
	// TODO: 예외처리
}

func GetSupporting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invitationCode := vars["invitationCode"]
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	supportings := supportingMap[invitationCode]
	if supportings != nil {
		_ = json.NewEncoder(w).Encode(supportings)
	} else {
		_ = json.NewEncoder(w).Encode([]supporting{})
	}
}

func MakeSupporting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	uIntUserId, _ := strconv.ParseUint(userId, 10, 32)
	id := uint(uIntUserId)
	user := entity.User{Id: &id}
	database.Connector.First(&user)
	invitationCode := vars["invitationCode"]

	supportings := supportingMap[invitationCode]
	supportings = append(supportings, supporting{UserName: user.UserName, CreatedTime: entity.DateTime{Time: time.Now()}})
	supportingMap[invitationCode] = supportings

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func SyncTimeRecords(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := io.ReadAll(r.Body)
	var commands []SyncTimerRecordCommand
	_ = json.Unmarshal(requestBody, &commands)

	if len(commands) == 0 {
		return
	}

	vars := mux.Vars(r)
	userId := vars["id"]
	uIntUserId, _ := strconv.ParseUint(userId, 10, 32)
	id := uint(uIntUserId)

	var timeRecords []entity.TimeRecord

	err := database.Connector.Transaction(func(tx *gorm.DB) error {
		user := entity.User{Id: &id}
		database.Connector.First(&user)

		timerStatesBylocalTimerId := make(map[uint]TimerState)
		for _, command := range commands {
			timerState := timerStatesBylocalTimerId[command.TimerId]
			timerState.TimerName = command.TimerName
			if timerState.StartTime == nil {
				timerState.StartTime = &DateTime{Time: command.StartTime.Time}
			} else {
				minTime := (*timerState.StartTime).Min(command.StartTime)
				timerState.StartTime = &minTime
			}
			timerState.EndTime = timerState.EndTime.Max(command.EndTime)
			timerStatesBylocalTimerId[command.TimerId] = timerState
		}

		timerIdsByLocalTimerId := make(map[uint]uint)
		for timerId, timerState := range timerStatesBylocalTimerId {
			var count int64
			tx.Model(&entity.Timer{}).Where("start_time = ?", *timerState.StartTime).Count(&count)

			if count == 0 {
				timer := entity.Timer{
					Name:      timerState.TimerName,
					MakerId:   id,
					Type:      entity.Personal,
					StartTime: *timerState.StartTime,
					EndTime:   &timerState.EndTime,
				}

				if err := tx.Create(&timer).Error; err != nil {
					return err
				}
				timerIdsByLocalTimerId[timerId] = *timer.Id
			}
		}

		for _, timeRecord := range commands {
			var count int64
			tx.Model(&entity.TimeRecord{}).Where("start_time = ?", timeRecord.StartTime).Count(&count)

			if count == 0 {
				timeRecords = append(timeRecords, entity.TimeRecord{TimerId: timerIdsByLocalTimerId[timeRecord.TimerId], UserId: id, StartTime: timeRecord.StartTime, EndTime: timeRecord.EndTime})
			}
		}

		if len(timeRecords) == 0 {
			return errors.New("duplicated time records")
		}

		if err := tx.CreateInBatches(&timeRecords, len(timeRecords)).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(timeRecords)
	}
}

type TimerState struct {
	TimerName string    `json:"timerName"`
	StartTime *DateTime `json:"startTime"`
	EndTime   DateTime  `json:"endTime"`
}

type SyncTimerRecordCommand struct {
	TimerName string   `json:"timerName"`
	TimerId   uint     `json:"timerId"`
	StartTime DateTime `json:"startTime"`
	EndTime   DateTime `json:"endTime"`
}

type Participant struct {
	UserName string        `json:"userName"`
	ToDos    []entity.ToDo `json:"toDos"`
}

type supporting struct {
	UserName    string   `json:"userName"`
	CreatedTime DateTime `json:"createdTime"`
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
