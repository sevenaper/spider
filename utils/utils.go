package utils

import (
	"crypto/md5"
	"fmt"
	"spider/consts"
	"spider/model"
	"strconv"
	"time"
)

//GetCommentURL 获取评论URL
func GetCommentURL(params *model.CommentParams) string {
	url := consts.CommentUrlPrefix + "&id=" + params.AppID +
		"&startIndex=" + strconv.Itoa(params.StartIndex) + "&endIndex=" + strconv.Itoa(params.EndIndex)
	return url
}

//GetVersionURL 获取版本号URL
func GetVersionURL(params *model.VersionParams) string {
	url := consts.VersionUrlPrefix +
		"page=" + strconv.Itoa(params.Page) + "/" +
		"id=" + params.AppID + consts.VersionUrlSuffix
	return url
}

//GetVersionPages 获取版本号页数
func GetVersionPages(count int) int {
	pages := 0
	if count%50 == 0 {
		pages = count / 50
	} else {
		pages = count/50 + 1
	}
	return pages
}

//ConvertGoTimeToStd 将golang的日期转换为标准日期格式，不带T那种
func ConvertGoTimeToStd(src string) string {
	const GoStr = "2006-01-02T15:04:05+08:00"
	t, _ := time.Parse(GoStr, src)
	return t.Format(consts.TimeStr)
}

//FillLastCrawlTime 填充最后爬取时间
func FillLastCrawlTime() string {
	curTime := time.Now()
	duration, _ := time.ParseDuration("-1h")
	eTime := curTime.Add(duration)
	return eTime.Format(consts.TimeStr)
}

//GenMainKey 根据不同appId和commentId规则生成全局唯一ID
func GenMainKey(key string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(key)))
}