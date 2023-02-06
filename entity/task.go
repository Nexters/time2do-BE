package entity

type Task struct {
	ID              uint   `json:"id"`
	UserID          string `json:"userId"`
	Lock            int    `json:"lock"`
	Text            string `json:"text"`
	CreatedTime     string `json:"createdTime"`
	CompletedTime   string `json:"completedTime"`
	CompletedStatus string `json:"completedStatus"`
	DeletedTime     string `json:"deletedTime"`
	Modified_time   string `json:"modifiedTime"`
}
