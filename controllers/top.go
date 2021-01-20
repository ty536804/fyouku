package controllers

import (
	"fyoukuApi/models"
	"github.com/astaxie/beego"
)

type TopController struct {
	beego.Controller
}

// @Summer 根据频道获取排行榜
func (t *TopController) ChannelTop() {
	channelId, _ := t.GetInt("channelId")
	if channelId == 0 {
		t.Data["json"] = ReturnError(4001, "必须指定频道ID")
		t.ServeJSON()
	}
	num, videos, err := models.GetChannelTop(channelId)
	if err == nil {
		t.Data["json"] = ReturnSuccess(200, "success", videos, num)
	} else {
		t.Data["json"] = ReturnError(400, "没有相关内容")
	}
	t.ServeJSON()
}

// @Summer 根据类型获取榜单
func (t *TopController) TypeTop() {
	typeId, _ := t.GetInt("topId")
	if typeId == 0 {
		t.Data["json"] = ReturnError(4001, "必须指定类型")
		t.ServeJSON()
	}
	num, videos, err := models.GetTypeTop(typeId)
	if err == nil {
		t.Data["json"] = ReturnSuccess(200, "success", videos, num)
	} else {
		t.Data["json"] = ReturnError(400, "没有相关内容")
	}
	t.ServeJSON()
}
