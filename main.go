package main

import (
	"spider/consts"
)

func main() {
	S := NewAppleCommentSpider()
	G := NewCommentGraph()
	//定义任务列表，ID映射可查看consts.go
	K := []string{"抖音"}

	T := []string{}
	for _, v := range K {
		T = append(T, consts.GetAppMap()[v])
	}
	StartCrawl(S, G, T)

}