package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Tweet struct {
	Id         int              `json:"id" db:"tweet_id"`
	Content    string           `json:"content" db:"tweet_content"`
	CreatedAt  time.Time        `json:"createdAt" db:"tweet_created_at"`
	ModifiedAt pgtype.Timestamp `json:"modifiedAt" db:"tweet_modified_at"`
	LikeCount  int              `json:"likeCount" db:"tweet_like_count"`
	UserId     int              `json:"-" db:"tweet_user_id"`

	IsLiked bool `json:"isLiked" db:"tweet_is_liked"`
	User
}
