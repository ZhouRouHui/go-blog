package main

import (
	"database/sql"
	"fmt"
	"goblog/bootstrap"
	"goblog/pkg/database"
	"goblog/pkg/logger"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	"unicode/utf8"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// 全局 router 对象
var router *mux.Router

// 全局数据库对象
var db *sql.DB

// Article 文章
type Article struct {
	Title, Body string
	ID          int64
}

// Delete 删除文章
func (a Article) Delete() (int64, error) {
	rs, err := db.Exec("DELETE FROM articles WHERE id = " + strconv.FormatInt(a.ID, 10))
	if err != nil {
		return 0, nil
	}

	// 更新成功，跳转到文章详情页
	if n, _ := rs.RowsAffected(); n > 0 {
		return n, nil
	}

	return 0, nil
}

// getArticleByID 通过 id 获取文章
func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)

	return article, err
}

// getRouteVariable 获取路由中的参数
func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r) // Vars() 返回当前请求的路由变量
	return vars[parameterName]
}

// articlesEditHandler 文章编辑页面
func articlesEditHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 url 参数
	id := getRouteVariable("id", r)

	// 2. 读取对应的文章数据
	article, err := getArticleByID(id)

	// 3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 读取成功，显示表单
		updateURL, _ := router.Get("articles.update").URL("id", id)
		data := ArticlesFormData{
			Title:  article.Title,
			Body:   article.Body,
			URL:    updateURL,
			Errors: nil,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		logger.LogError(err)

		tmpl.Execute(w, data)
	}
}

// articlesUpdateHandler 更新文章
func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取路由参数
	id := getRouteVariable("id", r)

	// 2. 获取模型数据
	article, err := getArticleByID(id)

	// 3. 错误判断
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 not found")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 internal server error")
		}
	} else {
		// 请求数据校验
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		errors := validateArticleFormData(title, body)
		if len(errors) > 0 {
			// 数据校验有错误
			editURL, _ := router.Get("articles.edit").URL("id", id)
			data := ArticlesFormData{
				Title:  article.Title,
				Body:   article.Body,
				URL:    editURL,
				Errors: errors,
			}
			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			logger.LogError(err)

			tmpl.Execute(w, data)
		} else {
			// 数据校验成功，进行更新
			query := "UPDATE articles SET title = ?, body = ? WHERE id = ?"
			rs, err := db.Exec(query, title, body, id)

			if err != nil {
				logger.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 internal server error")
			}

			// 更新成功
			if n, _ := rs.RowsAffected(); n > 0 {
				showURL, _ := router.Get("articles.show").URL("id", id)
				http.Redirect(w, r, showURL.String(), http.StatusFound)
			} else {
				fmt.Fprint(w, "您没有做任何修改!")
			}
		}
	}
}

// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

// validateArticleFormData 文章参数校验
func validateArticleFormData(title, body string) map[string]string {
	errors := make(map[string]string)

	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	return errors
}

// articlesStoreHandler 创建文章
func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)
	// 检查是否有错误
	if len(errors) == 0 {
		lastInsertID, err := saveArticleToDB(title, body)
		if lastInsertID > 0 {
			fmt.Fprint(w, "插入成功，ID 为"+strconv.FormatInt(lastInsertID, 10))
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		storeURL, _ := router.Get("articles.store").URL()
		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}

		tmpl.Execute(w, data)
	}
}

// saveArticleToDB 将数据写入 articles 表中
func saveArticleToDB(title string, body string) (int64, error) {
	// 变量初始化
	var (
		id   int64
		err  error
		rs   sql.Result
		stmt *sql.Stmt
	)

	// 1. 获取一个 prepare 声明语句
	stmt, err = db.Prepare("INSERT INTO articles (title, body) VALUES(?,?)")
	if err != nil {
		return 0, err
	}

	// 2. 在此函数运行结束后关闭此语句，防止占用 SQL 连接
	defer stmt.Close()

	// 3. 执行请求，传参进入绑定的内容
	rs, err = stmt.Exec(title, body)
	if err != nil {
		return 0, err
	}

	// 4. 插入成功的话，会返回自增ID
	if id, err = rs.LastInsertId(); id > 0 {
		return id, nil
	}

	return 0, err
}

// forceHTMLMiddleware 添加返回头标识中间件
func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1.设置标头
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 2.继续处理请求
		next.ServeHTTP(w, r)
	})
}

// removeTrailingSlash 删除路由后面的斜杠
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

// articlesCreateHandler 博文表单页面
func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {

	storeURL, _ := router.Get("articles.store").URL()
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}
	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, data)
}

// articlesDeleteHandler 删除文章
func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取文章id
	id := getRouteVariable("id", r)

	// 2. 读取对应的文章数据
	article, err := getArticleByID(id)

	// 3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 not found")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 internal server error")
		}
	} else {
		// 4. 未出现问题，执行删除
		rowsAffected, err := article.Delete()

		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 internal server err")
		} else {
			if rowsAffected > 0 {
				// 重定向到文章列表页
				indexURL, _ := router.Get("articles.index").URL()
				http.Redirect(w, r, indexURL.String(), http.StatusFound)
			} else {
				// Edge case
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 not found")
			}
		}
	}
}

func main() {
	// 初始化数据库连接
	database.Initialize()
	db = database.DB
	bootstrap.SetupDB()

	// 路由初始化
	router = bootstrap.SetupRoute()

	// router := mux.NewRouter().StrictSlash(true)
	router.StrictSlash(true)

	// 文章创建
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	// 博客创建页面
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")
	// 博客编辑页面
	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
	// 编辑保存
	router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")
	// 删除文章
	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

	// 中间件：强制内容类型为 HTML
	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
