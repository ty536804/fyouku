package main

import (
	"encoding/json"
	"fmt"
	"fyoukuApi/Services/mq"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	beego.LoadAppConfig("ini", "../../conf/app.conf")
	defaultDb := beego.AppConfig.String("defaultdb")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", defaultDb, 30, 30)

	mq.ConsumerDlx("fyouku.commant.count", "fyouku_commant_count",
		"fyouku.commant.count.dlx", "fyouku_commant_count_dlx", 10000, callback)
}

func callback(s string) {
	type Data struct {
		VideoId    int
		EpisodesId int
	}
	var data Data
	err := json.Unmarshal([]byte(s), &data)
	if err == nil {
		o := orm.NewOrm()
		//修改视频的总评论数
		o.Raw("UPDATE video SET comment=comment+1 WHERE id=?", data.VideoId).Exec()
		//修改视频剧集的评论数
		o.Raw("UPDATE video_episodes SET comment=comment+1 WHERE id=?", data.EpisodesId).Exec()
		//创建一个简单的模式的MQ 把要传递的数据转换为json字符串
		videoObj := map[string]int{"VideoId": data.VideoId}
		videoJson, _ := json.Marshal(videoObj)
		mq.Publish("", "fyouku_top", string(videoJson))
	}
	fmt.Print("msg is:%s\n", s)
}
