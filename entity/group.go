package entity

import "time"

type Group struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Maker        string    `json:"makerId"`
	Participants string    `json:"participants`
	Tags         string    `json:"tags"`
	SetTime      time.Time `json:"setTime"`
}
