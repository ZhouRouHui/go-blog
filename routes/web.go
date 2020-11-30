package routes

import (
	"goblog/app/http/controller"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterWebRoutes 注册路由
func RegisterWebRoutes(r *mux.Router) {
	// 静态页面
	pc := new(controller.PagesController)
	// 首页
	r.HandleFunc("/", pc.Home).Methods("GET").Name("home")
	// 关于我们
	r.HandleFunc("/about", pc.About).Methods("GET").Name("about")
	// 404
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)

	// 文章相关页面
	ac := new(controller.ArticlesController)
	// 文章详情
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("articles.show")
	// 文章列表
	r.HandleFunc("/articles", ac.Index).Methods("GET").Name("articles.index")
	// 文章创建
	r.HandleFunc("/articles", ac.Store).Methods("POST").Name("articles.store")
	// 博客创建页面
	r.HandleFunc("/articles/create", ac.Create).Methods("GET").Name("articles.create")
	// 博客编辑页面
	r.HandleFunc("/articles/{id:[0-9]+}/edit", ac.Edit).Methods("GET").Name("articles.edit")
	// 编辑保存
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Update).Methods("POST").Name("articles.update")
	// 删除文章
	r.HandleFunc("/articles/{id:[0-9]+}/delete", ac.Delete).Methods("POST").Name("articles.delete")
}
