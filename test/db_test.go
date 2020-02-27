package test

import (
	"fmt"
	"spider/service/database"
	"testing"
)

func TestDB(t *testing.T) {
	database.InitDB()
	lastTime := "2019-02-01T16:23:05"
	rows, _ := database.QueryTasks(lastTime)
	for _, v := range rows.Rows {
		fmt.Printf("%+v\n", v)
	}
}
