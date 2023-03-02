package entity

type TimerType int

const (
	_ TimerType = iota // To prevent mapping zero value
	Personal
	Group
)

type ToDo struct {
	Id            uint      `json:"id"`
	UserId        uint      `gorm:"not null" json:"userId"`
	Content       string    `gorm:"not null" json:"content"`
	Completed     bool      `json:"completed"`
	CreatedTime   DateTime  `json:"createdTime"`
	CompletedTime *DateTime `json:"completedTime"`
}

type Timer struct {
	Id             *uint     `json:"id"`
	Name           string    `gorm:"not null" json:"name"`
	MakerId        uint      `gorm:"not null" json:"makerId"`
	Type           TimerType `gorm:"not null" json:"type"`
	Tag            *string   `json:"tag"`
	InvitationCode *string   `json:"invitationCode"`
	Users          []*User   `gorm:"many2many:participants;" json:"users"`
	StartTime      DateTime  `gorm:"not null" json:"startTime"`
	EndTime        *DateTime `json:"endTime"`
}

type TimeRecord struct {
	Id        uint     `json:"id"`
	UserId    uint     `json:"userId"`
	TimerId   uint     `json:"timerId"`
	StartTime DateTime `json:"startTime"`
	EndTime   DateTime `json:"endTime"`
}

type User struct {
	Id         *uint   `json:"id"`
	UserName   string  `json:"userName"`
	Password   string  `json:"password"`
	Onboarding bool    `json:"onboarding"`
	Timers     []Timer `gorm:"many2many:participants;" json:"timers"`
}
