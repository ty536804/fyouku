package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"
)

//视频文件名生成函数
func GetVideoName(uid string) string {
	h := md5.New()
	h.Write([]byte(uid + strconv.FormatInt(time.Now().UnixNano(), 10)))
	return hex.EncodeToString(h.Sum(nil))
}
