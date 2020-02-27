package main

import (

	"spider/crons"
	"spider/service"
	"spider/service/database"
)

func Init() {
	database.InitDB()
	service.InitTaskService()
	crons.InitCrons()
}

func main() {
	S := service.NewAppleCommentSpider()
	G := service.NewCommentGraph()

	//任务列表
	K := service.GlobalTaskLoader.GetTaskMap()
	service.StartCrawl(S, G, K)

	crons.CronJobs()
}