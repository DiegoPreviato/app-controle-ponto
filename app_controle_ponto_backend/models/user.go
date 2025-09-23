package models

type User struct {
	ID       int64  `json:"id"`
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"` // omitempty so it's not sent in responses
}
