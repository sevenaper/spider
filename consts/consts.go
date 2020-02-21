package consts
//爬虫URL
const (
	COMMENT_URL_PREFIX = "https://itunes.apple.com/WebObjects/MZStore.woa/wa/userReviewsRow?cc=cn&displayable-kind=11&sort=0&appVersion=all"
	VERSION_URL_PREFIX = "https://itunes.apple.com/rss/customerreviews/"
	VERSION_URL_SUFFIX = "/sortby=mostrecent/json?l=en&&cc=cn"
)

//时间格式化模板
const TIME_STR = "2006-01-02 15:04:05"
var appMap = map[string]string{
	"抖音": "1142110895",
}

func GetAppMap() map[string]string {
	return appMap
}