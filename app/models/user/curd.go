package user

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/password"
	"goblog/pkg/types"
)

// Create 创建用户，通过 User.ID 来判断是否创建成功
func (u *User) Create() (err error) {
	if err = model.DB.Create(&u).Error; err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}

// Get 获取用户
func Get(uid string) (user User, err error) {
	id := types.StringToInt(uid)
	if err = model.DB.First(&user, id).Error; err != nil {
		return user, err
	}
	return user, nil
}

// GetByEmail 通过 email 获取用户
func GetByEmail(email string) (user User, err error) {
	if err = model.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

// ComparePassword 比对密码是否正确
func (u User) ComparePassword(_password string) bool {
	return password.CheckHash(_password, u.Password)
}
