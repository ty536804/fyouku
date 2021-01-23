package controllers

import (
	"fyoukuApi/Services/mq"
	"github.com/astaxie/beego"
	"strconv"
	"time"
)

type MqDemoController struct {
	beego.Controller
}

//简单模式和work工作模式 push方法
func (m *MqDemoController) GetMq() {
	go func() {
		count := 0
		for {
			mq.Publish("", "fyouku_demo", "hello"+strconv.Itoa(count))
			count++
			time.Sleep(1 * time.Second)
		}
	}()

	m.Ctx.WriteString("hello")
}
