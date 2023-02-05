package entity

type Group struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	MakerID      uint   `json:"makerId"`
	Participants string `json:"participants`
	Tags         string `json:"tags"`
	SetTime      string `json:"setTime"`
}
