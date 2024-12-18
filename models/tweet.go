package models

import "time"

type Tweet struct {
	Id         int       `json:"id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
	UserId     int       `json:"userId"`
}
