package database

import (
	"code.byted.org/gopkg/logs"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"spider/exceptions"
	"spider/model"
	"spider/utils"
)

var (
	dbReader *gorm.DB
)

func InitDB() {
	db, err := gorm.Open("mysql", "root:root9898008@tcp(134.175.33.133:3306)/gowatcher?charset=utf8&parseTime=True&loc=Local")
	if err == nil {
		db.SingularTable(true)
		dbReader = db
	} else {
		panic(err)
	}
}

//QueryTasks 获取爬虫任务列表
func QueryTasks(lastTime string) (*model.TaskTable, error) {
	fmt.Println(dbReader)
	rows, err := dbReader.Table("gowatcher.crawl_task_table").Debug().
		Select("id, app_name, app_id, status, create_time, modify_time").
		Where("modify_time > ?", lastTime).Order("modify_time").Rows()

	defer func() {
		if r := recover(); r != nil {
			logs.Error("Recovered in QueryTasks: %v", r)
		}

		//官方文档示例未考虑指针为空调用Close会panic的情形
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		logs.Warn("QueryTasks error, err: %v", err.Error())
		return nil, exceptions.ErrDBHandle
	}

	tTable := &model.TaskTable{}
	for rows.Next() {
		var tmpFT model.TaskRow
		dbReader.ScanRows(rows, &tmpFT)
		tTable.Rows = append(tTable.Rows, &tmpFT)
		tTable.LastTime = utils.ConvertGoTimeToStd(tmpFT.ModifyTime)
	}
	return tTable, nil
}
