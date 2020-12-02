package controller

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"goblog/pkg/view"
	"net/http"
	"unicode/utf8"

	"gorm.io/gorm"
)

// ArticlesController 文章相关页面
type ArticlesController struct{}

// Show 详情页面
func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 url 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)

	// 3. 如果出现问题
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
		// 4. 读取成功
		view.Render(w, _article, "articles.show")
	}
}

// Index 列表
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	// 1. 获取结果集
	articles, err := article.GetAll()
	if err != nil {
		// 数据库错误
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 internal server error")
	} else {
		// 2. 加载模板
		view.Render(w, articles, "articles.index")
	}
}

// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	Article     article.Article
	Errors      map[string]string
}

// Create 创建
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w, ArticlesFormData{}, "articles.create", "articles._form_field")
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

// Store 创建保存
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)
	// 检查是否有错误
	if len(errors) == 0 {
		_article := article.Article{
			Title: title,
			Body:  body,
		}
		_article.Create()
		if _article.ID > 0 {
			fmt.Fprint(w, "插入成功，ID 为"+types.Uint64ToString(_article.ID))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章失败")
		}
	} else {
		view.Render(w, ArticlesFormData{
			Title:  title,
			Body:   body,
			Errors: errors,
		}, "articles.store", "articles._form_field")
	}
}

// Edit 编辑页面
func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 url 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
		view.Render(w, ArticlesFormData{
			Title:  _article.Title,
			Body:   _article.Body,
			Article:    _article,
			Errors: nil,
		}, "articles.edit", "articles._form_field")
	}
}

// Update 编辑保存
func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	// 1. 获取路由参数
	id := route.GetRouteVariable("id", r)

	// 2. 获取模型数据
	_article, err := article.Get(id)

	// 3. 错误判断
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
		if len(errors) == 0 {
			// 数据校验成功，进行更新
			_article.Title = title
			_article.Body = body

			rowsAffected, err := _article.Update()
			if err != nil {
				// 数据库错误
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 internal server error")
			}
			if rowsAffected > 0 {
				showURL := route.Name2URL("articles.show", "id", id)
				http.Redirect(w, r, showURL, http.StatusFound)
			} else {
				fmt.Fprint(w, "您没有做任何更改！")
			}
		} else {
			// 数据校验有错误
			view.Render(w, ArticlesFormData{
				Title:  title,
				Body:   body,
				Article:    _article,
				Errors: errors,
			})
		}
	}
}

// Delete 删除文章
func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	// 1. 获取文章id
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 not found")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 internal server error")
		}
	} else {
		// 4. 未出现问题，执行删除
		rowsAffected, err := _article.Delete()

		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 internal server err")
		} else {
			if rowsAffected > 0 {
				// 重定向到文章列表页
				indexURL := route.Name2URL("articles.index")
				http.Redirect(w, r, indexURL, http.StatusFound)
			} else {
				// Edge case
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 not found")
			}
		}
	}
}
