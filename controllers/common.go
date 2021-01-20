package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"time"
)

type CommonController struct {
	beego.Controller
}

type JsonStruct struct {
	Code  int         `json:"code"`
	Msg   interface{} `json:"msg"`
	Items interface{} `json:"items"`
	Count int64       `json:"count"`
}

func ReturnSuccess(code int, msg, items interface{}, count int64) *JsonStruct {
	return &JsonStruct{
		Code:  code,
		Msg:   msg,
		Items: items,
		Count: count,
	}
}

func ReturnError(code int, msg interface{}) *JsonStruct {
	return &JsonStruct{
		Code: code,
		Msg:  msg,
	}
}

// 加密
func MD5V(password string) string {
	h := md5.New()
	h.Write([]byte(password + beego.AppConfig.String("md5code")))
	return hex.EncodeToString(h.Sum(nil))
}

// 格式化时间
func DateFormat(times int64) string {
	video_time := time.Unix(times, 0)
	return video_time.Format("2006-01-02")
}
