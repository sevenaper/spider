package model

//TaskRow 数据库行结构体
type TaskRow struct {
	ID         int32  `json:"id" gorm:"column:id"`
	AppID      string `json:"app_id" gorm:"column:app_id"`
	AppName    string `json:"app_name" gorm:"column:app_name"`
	Status     int32  `json:"status" gorm:"column:status"`
	CreateTime string `json:"create_time" gorm:"column:create_time"`
	ModifyTime string `json:"modify_time" gorm:"column:modify_time"`
}

//TaskTable 词表结构体
type TaskTable struct {
	Rows     []*TaskRow
	LastTime string
}
