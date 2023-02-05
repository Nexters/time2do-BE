package entity

type User struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
}
