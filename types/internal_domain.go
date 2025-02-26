package types

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	Id                 int
	Username           string
	FullName           string
	Email              string
	Password           string
	Token              string
	ProfileImage       pgtype.Text
	FollowerCount      int
	FollowingCount     int
	RecentTweetsLength int
	RecentTweets       []Tweet
}

type Tweet struct {
	Id        int
	Content   string
	CreatedAt time.Time
	LikeCount int
	User      User
}

type TweetWithUser struct {
	Tweet
	User User
}
