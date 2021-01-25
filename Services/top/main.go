package main

import (
	"encoding/json"
	"fmt"
	"fyoukuApi/Services/mq"
	redisClient "fyoukuApi/Services/redis"
	"fyoukuApi/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

func main() {
	beego.LoadAppConfig("ini", "../../conf/app.conf")
	defaultDb := beego.AppConfig.String("defaultdb")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", defaultDb, 30, 30)

	mq.Consumer("", "fyouku_top", callback)
}

func callback(s string) {
	type Data struct {
		VideoId int
	}
	var data Data
	err := json.Unmarshal([]byte(s), &data)
	videoInfo, err := models.RedisGetVideoInfo(data.VideoId)
	if err == nil {
		conn := redisClient.PoolConnect()
		defer conn.Close()
		//更新排行榜
		redisChannelKey := "video:top:channel:channelId" + strconv.Itoa(videoInfo.Id)
		redisTypeKey := "video:top:typeId:" + strconv.Itoa(videoInfo.TypeId)
		conn.Do("zincrby", redisChannelKey, 1, data.VideoId)
		conn.Do("zincrby", redisTypeKey, 1, data.VideoId)
	}
	fmt.Print("msg is:%s\n", s)
}
