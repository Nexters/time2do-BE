package controller

import (
	"encoding/json"
	"io"
	"net/http"
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
func GetAllTimers(w http.ResponseWriter, r *http.Request) {
	var timers []entity.Timer
	database.Connector.Find(&timers)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(timers)
}
