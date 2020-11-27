package article

import (
	"goblog/pkg/route"
	"strconv"
)

// Article 文章模型
type Article struct {
	ID    int
	Title string
	Body  string
}

// Link 生成文章链接
func (a Article) Link() string {
	return route.Name2URL("articles.show", "id", strconv.Itoa(a.ID))
}
