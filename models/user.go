package models

type User struct {
	BaseModel
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}
