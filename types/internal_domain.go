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
	IsFollowed         bool
	NextPageId         *int
}

type Tweet struct {
	Id         int
	Content    string
	CreatedAt  time.Time
	ModifiedAt time.Time
	LikeCount  int
	User       User
	IsLiked    bool
}

type Feed struct {
	Tweets     []Tweet
	NextPageId *int
}

type TweetInteraction struct {
	TweetId int
	IsLiked bool
}
