package entity

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"firstName"`
	Password string `json:"password"`
}
