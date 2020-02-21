package model

type CommentSpider struct {
	CommentId   string `json:"userReviewId"`
	Title       string `json:"title"`
	Content     string `json:"body"`
	Rating      int64  `json:"rating"`
	PublishTime string `json:"date"`
}

type VersionSpider struct {
	CommentId string
	AppName   string
	Version   string
}

type Comment struct {
	CommentId   string
	MainId      string
	Content     string
	Rating      string
	Version     string
	PublishTime string
	CrawlTime   string
}
