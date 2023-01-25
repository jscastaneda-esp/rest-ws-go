package models

import "time"

type BaseModel struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}
