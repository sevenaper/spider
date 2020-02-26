package service

//任务加载器全局句柄
var (
	GlobalTaskLoader *TaskLoader
)

//InitTaskService 初始化任务服务
func InitTaskService() {
	GlobalTaskLoader = NewTaskLoader()
}