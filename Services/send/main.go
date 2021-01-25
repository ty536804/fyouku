package main

import (
	"encoding/json"
	"fmt"
	"fyoukuApi/Services/mq"
	"fyoukuApi/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	beego.LoadAppConfig("ini", "../../conf/app.conf")
	defaultDb := beego.AppConfig.String("defaultdb")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", defaultDb, 30, 30)

	mq.Consumer("", "fyouku_send_message_user", callback)
}

func callback(s string) {
	type Data struct {
		UserId    int
		MessageId int64
	}
	var data Data
	err := json.Unmarshal([]byte(s), &data)
	if err == nil {
		models.SendMessageUser(data.UserId, data.MessageId)
	}
	fmt.Printf("msg is:%s\n", s)
}
