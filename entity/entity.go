package entity

import "time"

type TimerType int

const (
	_ TimerType = iota // To prevent mapping zero value
	Personal
	Group
)

type ToDo struct {
	Id            uint      `json:"id"`
	UserId        uint      `json:"userId"`
	Content       string    `gorm:"not null" json:"content"`
	Completed     bool      `json:"completed"`
	Private       bool      `json:"private"`
	CreatedTime   time.Time `json:"createdTime"`
	CompletedTime time.Time `json:"completedTime"`
	ModifiedTime  time.Time `json:"modifiedTime"`
	DeletedTime   time.Time `json:"deletedTime"`
}

type Timer struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	MakerId   uint      `json:"makerId"`
	Type      TimerType `gorm:"not null" json:"type"`
	Tags      string    `json:"tags"`
	LinkUrl   *string   `json:"linkUrl"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

type User struct {
	Id       uint   `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Participant struct {
	Id        uint      `json:"id"`
	UserId    uint      `json:"userId"`
	TableId   uint      `json:"tableId"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}
