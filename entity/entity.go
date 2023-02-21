package entity

type TimerType int

const (
	_ TimerType = iota // To prevent mapping zero value
	Personal
	Group
)

type ToDo struct {
	Id            uint     `json:"id"`
	UserId        uint     `gorm:"not null" json:"userId"`
	Content       string   `gorm:"not null" json:"content"`
	Completed     bool     `json:"completed"`
	CreatedTime   DateTime `json:"createdTime"`
	CompletedTime DateTime `json:"completedTime"`
	ModifiedTime  DateTime `json:"modifiedTime"`
	DeletedTime   DateTime `json:"deletedTime"`
}

type Timer struct {
	Id        uint      `json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	MakerId   uint      `gorm:"not null" json:"makerId"`
	Type      TimerType `gorm:"not null" json:"type"`
	Tags      string    `json:"tags"`
	LinkUrl   *string   `json:"linkUrl"`
	Users     []*User   `gorm:"many2many:participants;"`
	StartTime DateTime  `gorm:"not null" json:"startTime"`
	EndTime   *DateTime `json:"endTime"`
}

type TimeRecord struct {
	Id        uint     `json:"id"`
	UserId    uint     `json:"userId"`
	TimerId   uint     `json:"timerId"`
	StartTime DateTime `json:"startTime"`
	EndTime   DateTime `json:"endTime"`
}

type User struct {
	Id         uint    `json:"id"`
	UserName   string  `json:"username"`
	Password   string  `json:"password"`
	Onboarding bool    `json:"onboarding"`
	Timers     []Timer `gorm:"many2many:participants;"`
}