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
	id := uint(uIntUserId)
	_ = database.Connector.
		Where(&entity.User{Id: &id}).
		Preload("Timers").
		Preload("Timers.Users", func(db *gorm.DB) *gorm.DB {
			return db.Select("users.id")
		}).
		Find(&user)

	totalDuration := time.Duration(0)
	for _, timeRecord := range timeRecords {
		totalDuration += timeRecord.EndTime.Sub(timeRecord.StartTime)
	}

	report := Report{
		UserName:             user.UserName,
		TimeBlocksByDateTime: toTimeBlocks(timeRecords, user.Timers, toDos, firstDayOfMonth, lastDayOfMonth),
		TotalDuration:        totalDuration.String(),
	}
	_ = json.NewEncoder(w).Encode(report)
}

type Report struct {
	UserName             string               `json:"userName"`
	TimeBlocksByDateTime map[string]TimeBlock `json:"timeBlocks"`
	TotalDuration        string               `json:"totalDuration"`
}

type TimeBlock struct {
	Hour         int           `json:"hour"`
	Minute       int           `json:"minute"`
	ToDos        []entity.ToDo `json:"toDos"`
	GroupTimers  []GroupTimer  `json:"groupTimers"`
	InGroupTimer bool          `json:"inGroupTimer"`
}

type GroupTimer struct {
	Name              string   `json:"name"`
	DisplayTime       DateTime `json:"displayTime"`
	Duration          int      `json:"duration"`
	ParticipantsCount int      `json:"participantsCount"`
	Tag               string   `json:"tag"`
}

func toTimeBlocks(timeRecords []entity.TimeRecord, totalGroupTimers []entity.Timer, totalToDos []entity.ToDo, startDateTime DateTime, endDate DateTime) map[string]TimeBlock {
	var nowDate = startDateTime
	var timeBlocksByDateTime = map[string]TimeBlock{}

	for nowDate.Before(endDate) {
		hour := 0
		minute := 0
		inGroupTimer := false
		nextDate := nowDate.AddDate(0, 0, 1)
		for _, timeRecord := range timeRecords {
			actualEndDate := nextDate.Min(timeRecord.EndTime)

			if timeRecord.StartTime.Between(nowDate, nextDate) {
				duration := actualEndDate.Sub(timeRecord.StartTime)
				totalMinutes := int(duration.Minutes())
				hour = totalMinutes / 60
				minute = totalMinutes % 60
			} else if timeRecord.EndTime.Between(nowDate, nextDate) {
				duration := actualEndDate.Sub(nowDate)
				totalMinutes := int(duration.Minutes())
				hour = totalMinutes / 60
				minute = totalMinutes % 60
			}
		}

		//goland:noinspection GoPreferNilSlice
		var todayGroupTimers = []GroupTimer{}

		for _, timer := range totalGroupTimers {
			if timer.Type != entity.Group || timer.EndTime == nil {
				continue
			}

			inGroupTimer = timer.EndTime.Between(nowDate, nextDate)
			if inGroupTimer {
				todayGroupTimers = append(todayGroupTimers, GroupTimer{
					Name:              timer.Name,
					DisplayTime:       *timer.EndTime,
					Duration:          int(timer.EndTime.Sub(timer.StartTime).Hours()),
					ParticipantsCount: len(timer.Users),
					Tag:               timer.Tags,
				})
			}
		}

		//goland:noinspection GoPreferNilSlice
		var todayToDos = []entity.ToDo{}

		for _, toDo := range totalToDos {
			if toDo.CompletedTime.Between(nowDate, nextDate) {
				todayToDos = append(todayToDos, toDo)
			}
		}

		timeBlocksByDateTime[nowDate.Format("2006-01-02")] = TimeBlock{Hour: hour, Minute: minute, InGroupTimer: inGroupTimer, ToDos: todayToDos, GroupTimers: todayGroupTimers}
		//goland:noinspection GoPreferNilSlice
		var nowToDos = []entity.ToDo{}

		for _, toDo := range totalToDos {
			if toDo.CompletedTime.Between(nowDate, nextDate) {
				nowToDos = append(nowToDos, toDo)
			}
		}

		nowDate = nowDate.AddDate(0, 0, 1)
	}
	return timeBlocksByDateTime
}

// TODO: 응원하기
