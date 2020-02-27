package main

import (

	"spider/crons"
	"spider/service"
)

func main() {
	S := service.NewAppleCommentSpider()
	G := service.NewCommentGraph()
	//定义任务列表，ID映射可查看consts.go
	K := service.GlobalTaskLoader.GetTaskMap()
	service.StartCrawl(S, G, K)
	crons.CronJobs()

}