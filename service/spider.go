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

//Graph 总存储结构
type Graph map[string]*model.Comment

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

//NewGraph 初始化存储对象
func NewGraph() Graph {
	g := make(Graph)
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
		comment.Rating = val.Get("rating").String()
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
		comment.Rating = val.Get("rating").String()
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
func CrawlComment(s *AppleCommentSpider, g CommentGraph, t *model.Task) int {
	tmp := make(CommentGraph)
	params := model.CommentParams{
		AppID:      t.AppID,
		StartIndex: 0,
		EndIndex:   consts.PageSize,
	}
	url := utils.GetCommentURL(&params)
	s.Crawl(url)
	recentTime, hit := s.ParseFirstCommentContent(tmp, t.LastCrawlTime)
	time.Sleep(1 * time.Second)

	//如果第一次就命中，则无需多页爬取
	for !hit {
		params.StartIndex += consts.PageSize
		params.EndIndex += consts.PageSize
		url = utils.GetCommentURL(&params)
		s.Crawl(url)
		hit = s.ParsePagesCommentContent(tmp, t.LastCrawlTime)
		time.Sleep(1 * time.Second)
	}

	if recentTime > t.LastCrawlTime {
		t.LastCrawlTime = recentTime
	}

	if len(tmp) > 0 {
		for k, v := range tmp {
			g[t.AppID+"|"+k] = v
		}
	}

	return len(tmp)
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

//Crawl 获取指定链接的内容
func (s *AppleVersionSpider) Crawl(url string) error {
	if err := s.Downloader.Visit(url); err != nil {
		log.Println("Visit:", url, " [error]", err)
		return err
	}
	s.Downloader.Wait()
	log.Println("Visit:", url, " [success]")
	return nil
}

//Result 获取响应字符串
func (s *AppleVersionSpider) Result() string {
	return s.Resp
}

//ParseVersionContent 解析版本号
func (s *AppleVersionSpider) ParseVersionContent(g VersionGraph) {
	base := s.Result()
	data := gjson.Get(base, "feed").Get("entry")
	data.ForEach(func(key, val gjson.Result) bool {
		commentId := val.Get("id").Get("label").String()
		version := val.Get("im:version").Get("label").String()
		g[commentId] = version
		return true
	})
}

//CrawlVersion 爬取版本号
func CrawlVersion(s *AppleVersionSpider, g VersionGraph, params model.VersionParams) {
	url := utils.GetVersionURL(&params)
	s.Crawl(url)
	s.ParseVersionContent(g)
	time.Sleep(1 * time.Second)
}

//Crawl 执行爬虫流程
func Crawl(k *AppleSpiders, g Graph, t *model.Task) {
	cg := make(CommentGraph)
	num := CrawlComment(k.AppleCommentSpider, cg, t)
	if num == 0 {
		return
	}

	vg := make(VersionGraph)
	pages := utils.GetVersionPages(num)
	params := model.VersionParams{}
	params.AppID = t.AppID
	for page := 1; page <= pages; page++ {
		params.Page = page
		CrawlVersion(k.AppleVersionSpider, vg, params)
	}

	for k, v := range cg {
		comment := &model.Comment{
			CommentId:   v.CommentId,
			MainId:      utils.GenMainKey(k),
			Content:     v.Title + "|" + v.Content,
			Rating:      v.Rating,
			Version:     "UNKNOWN",
			PublishTime: v.PublishTime,
			CrawlTime:   t.LastCrawlTime,
		}
		if _, ok := vg[v.CommentId]; ok {
			comment.Version = vg[v.CommentId]
		}
		g[k] = comment
	}
}

//StartCrawl 筛选爬虫任务
func StartCrawl(k *AppleSpiders, g Graph, tasks TaskDict) {
	t1 := &model.Task{
		AppID:         "1142110895",
		LastCrawlTime: "2019-11-09 18:55:07",
		Status:        consts.Normal,
	}
	t2 := &model.Task{
		AppID:         "1142110895",
		LastCrawlTime: "2019-11-09 18:55:07",
		Status:        consts.Normal,
	}
	tasks[t1.AppID] = t1
	tasks[t2.AppID] = t2
	for _, t := range tasks {
		if t.Status == consts.Normal {
			Crawl(k, g, t)
		}
	}
}
