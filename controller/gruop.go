package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time2do/database"
	"time2do/entity"
)

// @Summary 그룹 생성하기
// @Tags Group (다운 타이머)
// @Accept  json
// @Produce  json
// @Router /group [post]
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var group entity.Timer
	json.Unmarshal(requestBody, &group)
	database.Connector.Create(group)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}

// @Summary 그룹 조회하기
// @Tags Group (다운 타이머)
// @Accept  json
// @Produce  json
// @Router /groups [get]
func GetAllGroup(w http.ResponseWriter, r *http.Request) {
	var timers []entity.Timer
	database.Connector.Find(&timers)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(timers)
}
