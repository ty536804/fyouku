package controllers

import (
	"fyoukuApi/models"
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

// 获取频道地区列表
func (b *BaseController) ChannelRegion() {
	channelId, _ := b.GetInt("channelId")
	if channelId == 0 {
		b.Data["json"] = ReturnError(4001, "必须指定频道")
		b.ServeJSON()
	}
	num, regions, err := models.GetChannelRegion(channelId)
	if err == nil {
		b.Data["json"] = ReturnSuccess(200, "success", regions, num)
	} else {
		b.Data["json"] = ReturnError(4004, "没有相关内容")
	}
	b.ServeJSON()
}

// 获取频道类型列表
func (b *BaseController) ChannelType() {
	channelId, _ := b.GetInt("channelId")
	if channelId == 0 {
		b.Data["json"] = ReturnError(4001, "必须指定频道")
		b.ServeJSON()
	}
	num, regions, err := models.GetChannelType(channelId)
	if err == nil {
		b.Data["json"] = ReturnSuccess(200, "success", regions, num)
	} else {
		b.Data["json"] = ReturnError(4004, "没有相关内容")
	}
	b.ServeJSON()
}

// 根据传入参数获取视频列表
func (b *BaseController) ChannelVideo() {
	//频道ID
	channelId, _ := b.GetInt("channelId")
	if channelId == 0 {
		b.Data["json"] = ReturnError(4001, "必须指定频道")
		b.ServeJSON()
	}
	//获取频道地区ID
	regionId, _ := b.GetInt("regionId")
	//获取频道类型ID
	typeId, _ := b.GetInt("typeId")
	// 获取状态
	end := b.GetString("end")
	// 获取排序
	sort := b.GetString("sort")
	// 获取页码信息
	limit, _ := b.GetInt("limit")
	offset, _ := b.GetInt("offset")
	if limit == 0 {
		limit = 12
	}
	num, videos, err := models.GetChannelVideoList(channelId, regionId, typeId, limit, offset, end, sort)
	if err == nil {
		b.Data["json"] = ReturnSuccess(200, "success", videos, num)
	} else {
		b.Data["json"] = ReturnError(4004, "没有相关内容")
	}
	b.ServeJSON()
}

// @Summer 我的频道管理
func (b *BaseController) UserVideo() {
	uid, _ := b.GetInt("uid")
	if uid == 0 {
		b.Data["json"] = ReturnError(4001, "必须指定用户")
		b.ServeJSON()
	}
	num, videos, err := models.GetUserVideo(uid)
	if err == nil {
		b.Data["json"] = ReturnSuccess(200, "success", videos, num)
	} else {
		b.Data["json"] = ReturnError(4004, "没有相关内容")
	}
	b.ServeJSON()
}
