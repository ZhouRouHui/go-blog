package controllers

import (
	"goblog/pkg/view"
	"net/http"
)

type AuthController struct {}

// Register 注册页面
func (*AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}

// DoRegister 注册保存
func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {
	//
}
