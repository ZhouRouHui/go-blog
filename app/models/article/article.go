package article

import (
	"goblog/app/models"
	"goblog/pkg/route"
	"goblog/pkg/types"
)

// Article 文章模型
type Article struct {
	models.BaseModel

	Title string `gorm:"type:varchar(255);not null;default '';comment:标题" valid:"title"`
	Body  string `gorm:"type:varchar(255);not null;default '';comment:内容" valid:"body"`
}

// Link 生成文章链接
func (a Article) Link() string {
	return route.Name2URL("articles.show", "id", types.Uint64ToString(a.ID))
}
