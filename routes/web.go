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
}
