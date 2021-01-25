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

//订阅模式
func (m *MqDemoController) GetFanout() {
	go func() {
		count := 0
		for {
			mq.PublishEx("fyouku.demo.fanout", "fanout", "", "fanout"+strconv.Itoa(count))
			count++
			time.Sleep(1 * time.Second)
		}
	}()
	m.Ctx.WriteString("fanout")
}

//路由模式
func (m *MqDemoController) GetDirect() {
	go func() {
		count := 0
		for {
			if count%2 == 0 {
				mq.PublishEx("fyouku.demo.direct", "direct", "two", "directone"+strconv.Itoa(count))
			} else {
				mq.PublishEx("fyouku.demo.direct", "direct", "one", "directone"+strconv.Itoa(count))
			}
			count++
			time.Sleep(1 * time.Second)
		}
	}()
	m.Ctx.WriteString("direct")
}

//topic主题模式
func (m *MqDemoController) GetTopic() {
	go func() {
		count := 0
		for {
			for {
				if count%2 == 0 {
					mq.PublishEx("fyouku.demo.topic", "topic", "fyouku.video", "fyouku.video"+strconv.Itoa(count))
				} else {
					mq.PublishEx("fyouku.demo.topic", "topic", "user.fyouku", "user.fyouku"+strconv.Itoa(count))
				}
				count++
				time.Sleep(1 * time.Second)
			}
		}
	}()
	m.Ctx.WriteString("topicall")
}

func (m *MqDemoController) GetTopicTwo() {
	go func() {
		count := 0
		for {
			for {
				if count%2 == 0 {
					mq.PublishEx("fyouku.demo.topic", "topic", "a.frog.name", "a.frog.name"+strconv.Itoa(count))
				} else {
					mq.PublishEx("fyouku.demo.topic", "topic", "b.frog.uid", "b.frog.uid"+strconv.Itoa(count))
				}
				count++
				time.Sleep(1 * time.Second)
			}
		}
	}()
	m.Ctx.WriteString("topicTwo")
}

//死信队列 ttl过期 消息被拒绝 队列达到最大长度 都会造成死性队列
func (m *MqDemoController) GetDlx() {
	go func() {
		count := 0
		for {
			mq.PublishDlx("fyouku.dlx.a", "dlx"+strconv.Itoa(count))
			count++
			time.Sleep(1 * time.Second)
		}
	}()
	m.Ctx.WriteString("dlx")
}

func (m *MqDemoController) GetDlxTwo() {
	go func() {
		count := 0
		for {
			mq.PublishEx("fyouku.dlx.b", "fanout", "", "dlxtwo"+strconv.Itoa(count))
			count++
			time.Sleep(1 * time.Second)
		}
	}()
	m.Ctx.WriteString("dlxTwo")
}
