package entity

type Group struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	MakerID      uint   `json:"makerId"`
	Participants string `json:"participants`
	Tags         string `json:"tags"`
	StartTime    string `json:"startTime"`
	SetTime      string `json:"setTime"`
}
