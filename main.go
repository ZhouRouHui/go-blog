package main

import (
	_ "github.com/go-sql-driver/mysql"
	"goblog/app/http/middlewares"
	"goblog/bootstrap"
	"net/http"
)


func main() {
	// 初始化数据库连接
	bootstrap.SetupDB()

	// 路由初始化
	router := bootstrap.SetupRoute()

	http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
}
