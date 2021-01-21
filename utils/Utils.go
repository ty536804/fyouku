package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

//视频文件名生成函数
func GetVideoName(uid string) string {
	h := md5.New()
	h.Write([]byte(uid + strconv.FormatInt(time.Now().UnixNano(), 10)))
	return hex.EncodeToString(h.Sum(nil))
}

type ReturnSuccessJson struct {
	//必须的大写开头
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Items interface{} `json:"items"`
	Count int64       `json:"count"`
}

func ReturnSuccess(code int, msg string, items interface{}, count int64) string {
	jsonData := ReturnSuccessJson{Code: code, Msg: msg, Items: items, Count: count}
	if bytes, err := json.Marshal(jsonData); err == nil {
		return string(bytes)
	}
	return ""
}

type ReturnErrorJson struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func ReturnError(code int, err interface{}) string {
	var msg string
	switch err.(type) {
	case string:
		msg, _ = err.(string)
	default:
		msg = fmt.Sprintf("%s", err)
	}

	jsonData := ReturnErrorJson{Code: code, Msg: msg}
	if bytes, err := json.Marshal(jsonData); err == nil {
		return string(bytes)
	}
	return ""
}
