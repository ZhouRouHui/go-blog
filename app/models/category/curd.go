package category

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/route"
)

// Create 创建分类，通过 category.ID 来判断是否创建成功
func (c *Category) Create() (err error) {
	if err = model.DB.Create(&c).Error; err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}

// All 获取分类数据
func All() ([]Category, error) {
	var categories []Category
	if err := model.DB.Find(&categories).Error; err != nil {
		return categories, err
	}
	return categories, nil
}

// Link 方法用来生成文章链接
func (c Category) Link() string {
	return route.Name2URL("categories.show", "id", c.GetStringID())
}
