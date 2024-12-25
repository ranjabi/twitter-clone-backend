package models

type Feed struct {
	Tweets     []Tweet `json:"tweets"`
	NextPageId *int    `json:"nextPageId"`
}
