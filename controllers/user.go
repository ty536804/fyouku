package controllers

import (
	"fyoukuApi/models"
	"fyoukuApi/utils"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

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

// @Summer 发送内容
func (u *UserController) SendMessageDo() {
	uids := u.GetString("uids")
	content := u.GetString("content")
	if uids == "" {
		u.Data["json"] = ReturnError(4001, "请选择接收人")
		u.ServeJSON()
	}
	if content == "" {
		u.Data["json"] = ReturnError(4002, "请填写发送内容")
		u.ServeJSON()
	}
	messageId, err := models.SendMessageDo(content)
	if err == nil {
		uidConfig := strings.Split(uids, ",")
		for _, v := range uidConfig {
			userId, _ := strconv.Atoi(v)
			models.SendMessageUser(userId, messageId)
		}
		u.Data["json"] = ReturnSuccess(200, "success", 1, 1)
	} else {
		u.Data["json"] = ReturnError(4004, "发送失败，请联系客服")
	}
	u.ServeJSON()
}

// @Summer 上传视频
func (u *UserController) UploadVideo() {
	var (
		err   error
		title string
	)
	r := *u.Ctx.Request
	//获取表单提交的数据
	uid := r.FormValue("uid")
	//获取文件流
	file, header, _ := r.FormFile("file")
	//转换文件流为二进制
	b, _ := ioutil.ReadAll(file)
	//生成文件名
	filename := strings.Split(header.Filename, ".")
	filename[0] = utils.GetVideoName(uid)
	//文件保存的位置
	var fileDir = "/Users/bincao/Documents/goadmin/src/fyoukuApi/static/" + filename[0] + "." + filename[1]
	//播放地址
	var playUrl = "static/video/" + filename[0] + "." + filename[1]
	err = ioutil.WriteFile(fileDir, b, 0777)
	if err == nil {
		title = utils.ReturnSuccess(200, "success", playUrl, 1)
	} else {
		title = utils.ReturnError(5000, "上传失败,请联系客服")
	}
	u.Ctx.WriteString(title)
}
