package route

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router 路由对象
var Router *mux.Router

// Initialize 初始化路由
func Initialize() {
	Router = mux.NewRouter()
}

// Name2URL 路由名字转 url
func Name2URL(routeName string, pairs ...string) string {
	url, err := Router.Get(routeName).URL(pairs...)
	if err != nil {
		// checkError(err)
		return ""
	}

	return url.String()
}

// GetRouteVariable 获取路由中的参数
func GetRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r) // Vars() 返回当前请求的路由变量
	return vars[parameterName]
}
