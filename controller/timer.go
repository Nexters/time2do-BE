package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"time2do/database"
	"time2do/entity"
)

// @Summary 타이머 생성하기
// @Accept  json
// @Produce  json
// @Router /group [post]
func CreateTimer(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := io.ReadAll(r.Body)
	var timer entity.Timer
	_ = json.Unmarshal(requestBody, &timer)
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

type CreateTimeRecordCommand struct {
	StartTime entity.DateTime `json:"startDateTime"`
	EndTime   entity.DateTime `json:"endDateTime"`
}
