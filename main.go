package main

import (
	"fmt"
	"spider/crons"
	"spider/service"
)

func Init() {
	//database.InitDB()
	//service.InitTaskService()
	//crons.InitCrons()
}

func main() {
	S := service.NewAppleCommentSpider()
	G := service.NewCommentGraph()

	//任务列表
	//K := service.GlobalTaskLoader.GetTaskMap()
	//service.StartCrawl(S, G, K)
	service.StartCrawl(S, G, make(service.TaskDict))


	crons.CronJobs()
}