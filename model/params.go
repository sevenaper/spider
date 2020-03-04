package model

//CommentParams 评论请求参数
type CommentParams struct {
	AppID      string //苹果商店应用ID
	StartIndex int    //开始偏移量，最小0
	EndIndex   int    //结束偏移量，貌似无上限，推荐每次StartIndex+200
}

//VersionParams 版本号请求参数
type VersionParams struct {
	AppID string //苹果商店应用ID
	Page  int    //页数
}

//Task 爬取任务
type Task struct {
	AppID         string //苹果商店应用ID
	LastCrawlTime string //最后爬取时间
	Status        int32  //渠道是否启用
}