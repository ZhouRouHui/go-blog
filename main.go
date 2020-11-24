package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var router = mux.NewRouter()

// 首页
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog！</h1>")
}

// about 页面
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

// 404 not found 页面
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

// 文章详情
func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Vars() 返回当前请求的路由变量
	id := vars["id"]
	fmt.Fprint(w, "文章 ID："+id)
}

// 文章列表
func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "访问文章列表")
}

// 创建文章
func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "创建新的文章")
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1.设置标头
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 2.继续处理请求
		next.ServeHTTP(w, r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1.除首页之外，移除所有请求路径后面的斜杠
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		// 2. 将请求传下去
		next.ServeHTTP(w, r)
	})
}

// 博文表单页面
func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE hmtl>
<html lang="en">
<head>
	<title>创建文章 —— 我的技术博客</title>
</head>
<body>
	<form action="%s" method="post">
		<p><input type="text" name="title"></p>
		<p><textarea name="body" cols="30" rows="10"></textarea></p>
		<p><button type="submit">提交</button></p>
	</form>
</body>
</html>
`
	storeURL, _ := router.Get("articles.store").URL()
	fmt.Fprintf(w, html, storeURL)
}

func main() {
	// router := mux.NewRouter().StrictSlash(true)
	router.StrictSlash(true)

	// 首页
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	// 关于我们
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	// 文章详情
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	// 文章列表
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	// 文章创建
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	// 创建博文表单
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

	// 自定义 404 页面
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 中间件：强制内容类型为 HTML
	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
