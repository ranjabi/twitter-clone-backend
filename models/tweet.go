package models

import (
	"time"
)

type Tweet struct {
	Id               int        `json:"id"`
	Content          string     `json:"content"`
	CreatedAt        time.Time  `json:"createdAt" db:"created_at"`
	ModifiedAt       *time.Time `json:"modifiedAt" db:"modified_at"`
	UserId           int        `json:"userId" db:"user_id"`
	LikeCount        int        `json:"likeCount" db:"like_count"`
	IsLiked          bool       `json:"isLiked" db:"is_liked"`
	Username         string     `json:"username" db:"username"`
	UserFullName     string     `json:"userFullName" db:"full_name"`
	UserProfileImage string     `json:"userProfileImage" db:"profile_image"`
}
