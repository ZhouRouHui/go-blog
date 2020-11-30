package main

import (
	"goblog/app/http/middlewares"
	"goblog/bootstrap"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 初始化数据库连接
	// database.Initialize()
	// db = database.DB
	bootstrap.SetupDB()

	// 路由初始化
	router := bootstrap.SetupRoute()

	http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
}
