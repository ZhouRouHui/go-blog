package route

import (
	"goblog/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

var route *mux.Router

// SetRoute 设置路由实例，以供 Name2URL 等函数使用
func SetRoute(r *mux.Router) {
	route = r
}

// Name2URL 路由名字转 url
func Name2URL(routeName string, pairs ...string) string {
	url, err := route.Get(routeName).URL(pairs...)
	if err != nil {
		logger.LogError(err)
		return ""
	}

	return url.String()
}

// GetRouteVariable 获取路由中的参数
func GetRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r) // Vars() 返回当前请求的路由变量
	return vars[parameterName]
}
