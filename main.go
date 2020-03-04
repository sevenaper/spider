package main

import (
	"spider/service"
)

func Init() {
	//database.InitDB()
	//service.InitTaskService()
	//crons.InitCrons()
}

func main() {
	S := service.NewAppleSpiders()
	G := service.NewGraph()

	//任务列表
	//K := service.GlobalTaskLoader.GetTaskMap()
	//service.StartCrawl(S, G, K)
	service.StartCrawl(S, G, make(service.TaskDict))

	//crons.CronJobs()
}