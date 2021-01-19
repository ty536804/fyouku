package controllers

import (
	"fyoukuApi/models"
	"regexp"

	"github.com/astaxie/beego"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title register
// @Description 用户注册
// @Param	mobile			query 	string	true		"The mobile for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /register/login [post]
func (u *UserController) SaveRegister() {
	mobile := u.GetString("mobile")
	password := u.GetString("password")

	re := regexp.MustCompile(`^1\d{10}$`)
	if !re.MatchString(mobile) || len(mobile) < 11 {
		u.Data["json"] = ReturnError(4001, "请填写有效的手机号码")
		u.ServeJSON()
	}

	if len(password) < 1 {
		u.Data["json"] = ReturnError(4003, "密码不能为空")
		u.ServeJSON()
	}
	//判断手机号是否已经注册
	status := models.IsUserMobile(mobile)
	if status {
		u.Data["json"] = ReturnError(4005, "此手机号是否已经注册")
	} else {
		err := models.UserSave(mobile, MD5V(password))
		if err == nil {
			u.Data["json"] = ReturnSuccess(0, "注册成功", nil, 0)
			u.ServeJSON()
		} else {
			u.Data["json"] = ReturnError(5000, err)
		}
	}
	u.ServeJSON()
}

// @Title Login
// @Description 用户登录
// @Param	mobile		query 	string	true		"The mobile for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /register/login [post]
func (u *UserController) LoginDo() {
	mobile := u.GetString("mobile")
	password := u.GetString("password")

	re := regexp.MustCompile(`^1\d{10}$`)
	if !re.MatchString(mobile) || len(mobile) < 11 {
		u.Data["json"] = ReturnError(4001, "请填写有效的手机号码")
		u.ServeJSON()
	}

	if len(password) < 1 {
		u.Data["json"] = ReturnError(4003, "密码不能为空")
		u.ServeJSON()
	}
	uid, uname := models.IsMobileLogin(mobile, MD5V(password))
	if uid != 0 {
		u.Data["json"] = ReturnSuccess(200, "登录成功", map[string]interface{}{"uid": uid, "uname": uname}, 1)
	} else {
		u.Data["json"] = ReturnError(4004, "手机或密码不正确")
	}
	u.ServeJSON()
}
