package models

type Post struct {
	BaseModel
	PostContent string `json:"postContent"`
	UserId      string `json:"userId"`
}
