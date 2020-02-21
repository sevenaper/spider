package utils

import (
	"spider/consts"
	"spider/model"
	"strconv"
)

//GetCommentURL 获取评论URL
func GetCommentURL(key string, params *model.CommentParams) string {
	url := consts.COMMENT_URL_PREFIX + "&id=" + params.AppID +
		"&startIndex=" + strconv.Itoa(params.StartIndex) + "&endIndex=" + strconv.Itoa(params.EndIndex)
	return url
}

//GetVersionURL 获取版本号URL
func GetVersionURL(params *model.VersionParams) string {
	url := consts.VERSION_URL_PREFIX +
		"page=" + strconv.Itoa(params.Page) + "/" +
		"id=" + params.AppID + consts.VERSION_URL_SUFFIX
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
