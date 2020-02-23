package crons

import (
	"fmt"
	"github.com/robfig/cron"
)

//InitCrons 初始化定时更新任务
func InitCrons() {
	CronJobs()
}

//CronJobs 定时任务
func CronJobs() {
	i := 0
	c := cron.New()
	spec := "*/1 * * * *"
	c.AddFunc(spec, func() {
		i++
		fmt.Println("crons running:", i)
	})
	c.Start()

	select {}
}
