package user

import (
	"goblog/app/models"
	"goblog/pkg/route"
)

// User 用户模型
type User struct {
	models.BaseModel

	Name     string `gorm:"type:varchar(255);not null;unique;comment:昵称;" valid:"name"`
	Email    string `gorm:"type:varchar(255);unique;" valid:"email"`
	Password string `gorm:"type:varchar(255);" valid:"password"`
	// gorm:"-" -- 设置 gorm 在读写时略过此字段
	PasswordConfirm string `gorm:"-" valid:"password_confirm"`
}

// Link 用户链接
func (u User) Link() string {
	return route.Name2URL("users.show", "id", u.GetStringID())
}
