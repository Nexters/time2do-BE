package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
	database "time2do/database"
	"time2do/entity"
)

type DateTime = entity.DateTime

func ViewReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	params := r.URL.Query()
	parsed, _ := time.Parse("2006-01", params["yearMonth"][0])
	userId := vars["id"]
	uIntUserId, _ := strconv.ParseUint(userId, 10, 32)

	yearMonth := DateTime{Time: parsed}
	firstDayOfMonth := yearMonth.FirstDayOfMonth()
	lastDayOfMonth := yearMonth.LastDayOfMonth()

	// TODO: 트랜잭션 처리
	var timeRecords []entity.TimeRecord
	database.Connector.
		Where("user_id = ? AND start_time BETWEEN ? AND ?", userId, firstDayOfMonth, lastDayOfMonth).
		Find(&timeRecords)

	var toDos []entity.ToDo
	database.Connector.
		Where("user_id = ? AND completed = ? AND created_time BETWEEN ? AND ?", userId, true, firstDayOfMonth, lastDayOfMonth).
		Order("created_time desc").
		Find(&toDos)

	var user entity.User
	_ = database.Connector.
		Where(&entity.User{Id: uint(uIntUserId)}).
		Preload("Timers").
		Preload("Timers.Users", func(db *gorm.DB) *gorm.DB {
			return db.Select("users.id")
		}).
		Find(&user)

	var groupTimers []GroupTimer
	for _, timer := range user.Timers {
		if timer.Type != entity.Group || timer.EndTime == nil {
			continue
		}

		groupTimers = append(groupTimers, GroupTimer{
			Name:              timer.Name,
			DisplayTime:       *timer.EndTime,
			Duration:          int(timer.EndTime.Sub(timer.StartTime).Hours()),
			ParticipantsCount: len(timer.Users),
			Tag:               timer.Tags,
		})
	}

	totalDuration := time.Duration(0)
	for _, timeRecord := range timeRecords {
		totalDuration += timeRecord.EndTime.Sub(timeRecord.StartTime)
	}

	report := Report{
		ToDos:         toDos,
		TimeBlocks:    toTimeBlocks(timeRecords, user.Timers, firstDayOfMonth, lastDayOfMonth),
		GroupTimers:   groupTimers,
		TotalDuration: totalDuration,
	}
	_ = json.NewEncoder(w).Encode(report)
}

type Report struct {
	ToDos         []entity.ToDo `json:"toDos"`
	TimeBlocks    []TimeBlock   `json:"timeBlocks"`
	GroupTimers   []GroupTimer  `json:"groupTimers"`
	TotalDuration time.Duration `json:"duration"`
}

type TimeBlock struct {
	Hour         int  `json:"hour"`
	InGroupTimer bool `json:"inGroupTimer"`
}

type GroupTimer struct {
	Name              string   `json:"name"`
	DisplayTime       DateTime `json:"displayTime"`
	Duration          int      `json:"duration"`
	ParticipantsCount int      `json:"participantsCount"`
	Tag               string   `json:"tag"`
}

func toTimeBlocks(timeRecords []entity.TimeRecord, groupTimers []entity.Timer, startDateTime DateTime, endDate DateTime) []TimeBlock {
	var nowDate = startDateTime
	var timeBlocks []TimeBlock

	for nowDate.Before(endDate) {
		var hour = 0
		var inGroupTimer = false
		nextDate := nowDate.AddDate(0, 0, 1)
		for _, timeRecord := range timeRecords {
			actualEndDate := nextDate.Min(timeRecord.EndTime)

			if timeRecord.StartTime.Between(nowDate, nextDate) {
				hour = int(actualEndDate.Sub(timeRecord.StartTime).Hours())
			} else if timeRecord.EndTime.Between(nowDate, nextDate) {
				hour = int(actualEndDate.Sub(nowDate).Hours())
			}
		}

		for _, groupTimer := range groupTimers {
			inGroupTimer = groupTimer.EndTime != nil && groupTimer.EndTime.Between(nowDate, nextDate)
		}
		nowDate = nowDate.AddDate(0, 0, 1)
		timeBlocks = append(timeBlocks, TimeBlock{Hour: hour, InGroupTimer: inGroupTimer})
	}
	return timeBlocks
}

// TODO: 응원하기
