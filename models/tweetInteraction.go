package models

type TweetInteraction struct {
	TweetId int  `json:"tweetId" db:"tweet_id"`
	IsLiked bool `json:"isLiked" db:"is_liked"`
}
