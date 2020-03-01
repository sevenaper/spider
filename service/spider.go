package service

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/tidwall/gjson"
	"log"
	"spider/consts"
	"spider/model"
	"spider/utils"
	"time"
)

//CommentGraph 评论存储结构
type CommentGraph map[string]*model.CommentSpider

//VersionGraph 版本号存储结构
type VersionGraph map[string]string

//AppleCommentSpider 苹果应用商店评论爬虫
type AppleCommentSpider struct {
	Req        string
	Resp       string
	Err        string
	StatusCode int
	Downloader *colly.Collector
}

//AppleVersionSpider 苹果应用商店评论版本号爬虫
type AppleVersionSpider struct {
	Req        string
	Resp       string
	Err        string
	StatusCode int
	Downloader *colly.Collector
}

//AppleSpiders 苹果商店评论爬虫集合体
type AppleSpiders struct {
	AppleCommentSpider *AppleCommentSpider
	AppleVersionSpider *AppleVersionSpider
}

//NewCommentGraph 初始化评论存储对象
func NewCommentGraph() CommentGraph {
	g := make(CommentGraph)
	return g
}

//NewVersionGraph 初始化版本号存储对象
func NewVersionGraph() VersionGraph {
	g := make(VersionGraph)
	return g
}

//NewAppleCommentSpider 应用商店评论爬虫初始化
func NewAppleCommentSpider() *AppleCommentSpider {
	spider := &AppleCommentSpider{}
	spider.InitDownloader()
	return spider
}

//NewAppleVersionSpider 应用商店版本号爬虫初始化
func NewAppleVersionSpider() *AppleVersionSpider {
	spider := &AppleVersionSpider{}
	spider.InitDownloader()
	return spider
}

//NewAppleSpiders 爬虫初始化
func NewAppleSpiders() *AppleSpiders {
	spiders := &AppleSpiders{
		AppleCommentSpider: &AppleCommentSpider{},
		AppleVersionSpider: &AppleVersionSpider{},
	}
	spiders.AppleCommentSpider.InitDownloader()
	spiders.AppleVersionSpider.InitDownloader()
	return spiders
}

//InitDownloader 初始化评论下载器
func (s *AppleCommentSpider) InitDownloader() {
	s.Downloader = colly.NewCollector()
	s.Downloader.UserAgent = "iTunes/11.0 (Windows; Microsoft Windows 7 Business Edition Service Pack 1 (Build 7601)) AppleWebKit/536.27.1"
	s.Downloader.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", "apple.com")
		r.Headers.Set("Origin", "https://itunes.apple.com")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	})

	s.Downloader.OnResponse(func(r *colly.Response) {
		s.StatusCode = r.StatusCode
		s.Resp = string(r.Body)
		//fmt.Println(s.Resp, s.StatusCode)
	})

	s.Downloader.OnError(func(resp *colly.Response, errHttp error) {
		s.Err = errHttp.Error()
		fmt.Println(s.Err)
	})
}

//Result 获取响应字符串
func (s *AppleCommentSpider) Result() string {
	return s.Resp
}

//Crawl 获取指定链接的内容
func (s *AppleCommentSpider) Crawl(url string) error {
	if err := s.Downloader.Visit(url); err != nil {
		log.Println("Visit:", url, " [error]", err)
		return err
	}
	s.Downloader.Wait()
	log.Println("Visit:", url, " [success]")
	return nil
}

//ParseFirstCommentContent 解析第一次爬取的评论内容，并返回最近评论的时间
func (s *AppleCommentSpider) ParseFirstCommentContent(g CommentGraph, t string) (string, bool) {
	recentTime := t
	hit := false
	base := s.Result()
	data := gjson.Get(base, "userReviewList")
	data.ForEach(func(key, val gjson.Result) bool {
		comment := &model.CommentSpider{}
		comment.Title = val.Get("title").String()
		comment.Content = val.Get("body").String()
		comment.CommentId = val.Get("userReviewId").String()
		comment.PublishTime = val.Get("date").Time().In(time.Local).Format(consts.TimeStr)
		comment.Rating = val.Get("rating").Int()
		if comment.PublishTime > t {
			g[comment.CommentId] = comment
			if comment.PublishTime > recentTime {
				recentTime = comment.PublishTime
			}
		} else {
			hit = true
		}
		return true
	})

	return recentTime, hit
}

//ParsePagesCommentContent 解析爬取的评论内容，并返回是否爬取到上次时间点
func (s *AppleCommentSpider) ParsePagesCommentContent(g CommentGraph, t string) bool {
	hit := false
	base := s.Result()
	data := gjson.Get(base, "userReviewList")
	data.ForEach(func(key, val gjson.Result) bool {
		comment := &model.CommentSpider{}
		comment.Title = val.Get("title").String()
		comment.Content = val.Get("body").String()
		comment.CommentId = val.Get("userReviewId").String()
		comment.PublishTime = val.Get("date").Time().In(time.Local).Format(consts.TimeStr)
		comment.Rating = val.Get("rating").Int()
		if comment.PublishTime > t {
			g[comment.CommentId] = comment
		} else {
			hit = true
		}
		return true
	})
	return hit
}

//CrawlComment 爬取评论
func CrawlComment(s *AppleCommentSpider, g CommentGraph, t *model.Task) {
	params := model.CommentParams{
		AppID:      t.AppID,
		StartIndex: 0,
		EndIndex:   consts.PageSize,
	}
	url := utils.GetCommentURL(&params)
	s.Crawl(url)
	recentTime, hit := s.ParseFirstCommentContent(g, t.LastCrawlTime)

	//如果第一次就命中，则无需多页爬取
	for !hit {
		params.StartIndex += consts.PageSize
		params.EndIndex += consts.PageSize
		url = utils.GetCommentURL(&params)
		s.Crawl(url)
		hit = s.ParsePagesCommentContent(g, t.LastCrawlTime)
	}

	fmt.Println(len(g))

	if recentTime > t.LastCrawlTime {
		t.LastCrawlTime = recentTime
	}
	fmt.Println(t.LastCrawlTime)
}

//InitDownloader 初始化版本号下载器
func (s *AppleVersionSpider) InitDownloader() {
	s.Downloader = colly.NewCollector()
	s.Downloader.UserAgent = "iTunes/11.0 (Windows; Microsoft Windows 7 Business Edition Service Pack 1 (Build 7601)) AppleWebKit/536.27.1"
	s.Downloader.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", "itunes.apple.com")
		r.Headers.Set("Origin", "https://itunes.apple.com")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	})

	s.Downloader.OnResponse(func(r *colly.Response) {
		s.StatusCode = r.StatusCode
		s.Resp = string(r.Body)
		//fmt.Println(s.Resp, s.StatusCode)
	})

	s.Downloader.OnError(func(resp *colly.Response, errHttp error) {
		s.Err = errHttp.Error()
		fmt.Println(s.Err)
	})
}

//Crawl 执行爬取流程
func Crawl(k *AppleCommentSpider, g CommentGraph, t *model.Task) {
	CrawlComment(k, g, t)
}

//StartCrawl 筛选爬虫任务
func StartCrawl(k *AppleCommentSpider, g CommentGraph, tasks TaskDict) {
	//for _, t := range tasks {
	//	if t.Status == consts.Normal {
	//		Crawl(k, g, t)
	//	}
	//}
	t := &model.Task{
		AppID:         "1142110895",
		LastCrawlTime: "2019-11-09 16:55:07",
		Status:        consts.Normal,
	}
	CrawlComment(k, g, t)
