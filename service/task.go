package service

import (
	"spider/model"
	"spider/service/database"
	"spider/utils"
)

//TaskDict 存储app_id + 最后爬取时间 + 任务状态
type TaskDict map[string]*model.Task

//TaskLoader 任务加载器
type TaskLoader struct {
	TaskMap  TaskDict //任务列表
	LastTime string   //任务列表最近更新时间
}

//NewTaskLoader 创建任务加载器
func NewTaskLoader() *TaskLoader {
	return &TaskLoader{
		TaskMap:  make(TaskDict, 0),
		LastTime: "2019-01-01 00:00:00",
	}
}

//Load 加载任务
func (p *TaskLoader) Load() {
	rows, err := database.QueryTasks(p.LastTime)
	if err != nil {
		panic(err)
	}

	if rows.LastTime > p.LastTime {
		p.LastTime = rows.LastTime
		for _, v := range rows.Rows {
			p.ModifyTask(v)
		}
	}
}

//ModifyTask 更改任务列表
func (p *TaskLoader) ModifyTask(task *model.TaskRow) {
	if _, ok := p.TaskMap[task.AppID]; !ok {
		p.TaskMap[task.AppID] = &model.Task{
			AppID:         task.AppID,
			LastCrawlTime: utils.FillLastCrawlTime(),
			Status:        task.Status,
		}
	} else {
		p.TaskMap[task.AppID].Status = task.Status
	}
}

//GetTaskMap 获取任务列表
func (p *TaskLoader) GetTaskMap() TaskDict {
	return p.TaskMap
}