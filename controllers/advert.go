package controllers

import (
	"fyoukuApi/models"
	"github.com/astaxie/beego"
)

type AdvertController struct {
	beego.Controller
}

func (v *AdvertController) ChannelAdvert() {
	channelId, _ := v.GetInt("channelId")
	if channelId == 0 {
		v.Data["json"] = ReturnError(4001, "必须指定频道")
		v.ServeJSON()
	}
	num, videos, err := models.GetChannelAdvert(channelId)
	if err == nil {
		v.Data["json"] = ReturnSuccess(200, "success", videos, num)
	} else {
		v.Data["json"] = ReturnError(4004, "请求数据失败,请稍后重试~")
	}
	v.ServeJSON()
}

// @Summer 频道页-获取正在热播
func (v *AdvertController) ChanelHotList() {
	channelId, _ := v.GetInt("channelId")
	if channelId == 0 {
		v.Data["json"] = ReturnError(4001, "必须指定频道")
		v.ServeJSON()
	}
	num, videos, err := models.GetChannelHotList(channelId)
	if err == nil {
		v.Data["json"] = ReturnSuccess(200, "success", videos, num)
	} else {
		v.Data["json"] = ReturnError(4004, "没有相关内容")
	}
	v.ServeJSON()
}

// @Summer 频道页-根据频道地区获取推荐的视频
func (v *AdvertController) ChannelRecommendRegionList() {
	channelId, _ := v.GetInt("channelId")
	regionId, _ := v.GetInt("regionId")
	if channelId == 0 {
		v.Data["json"] = ReturnError(4001, "必须指定频道")
		v.ServeJSON()
	}
	if regionId == 0 {
		v.Data["json"] = ReturnError(4002, "必须指定频道地区")
		v.ServeJSON()
	}
	num, videos, err := models.GetChannelRecommendRegionList(channelId, regionId)
	if err == nil {
		v.Data["json"] = ReturnSuccess(200, "success", videos, num)
	} else {
		v.Data["json"] = ReturnError(4004, "没有相关内容")
	}
	v.ServeJSON()
}

// @Summer 频道页-根据频道类型获取推荐的视频
func (v *AdvertController) GetChannelRecommendTypeList() {
	channelId, _ := v.GetInt("channelId")
	typeId, _ := v.GetInt("typeId")
	if channelId == 0 {
		v.Data["json"] = ReturnError(4001, "必须指定频道")
		v.ServeJSON()
	}
	if typeId == 0 {
		v.Data["json"] = ReturnError(4002, "必须指定频道类型")
		v.ServeJSON()
	}
	num, videos, err := models.GetChannelRecommendTypeList(channelId, typeId)
	if err == nil {
		v.Data["json"] = ReturnSuccess(200, "success", videos, num)
	} else {
		v.Data["json"] = ReturnError(4004, "没有相关内容")
	}
	v.ServeJSON()
}

// @Summer 获取视频详情
func (v *AdvertController) VideoInfo() {
	videoId, _ := v.GetInt("videoId")
	if videoId == 0 {
		v.Data["json"] = ReturnError(4001, "必须指定视频ID")
		v.ServeJSON()
	}
	video, err := models.RedisGetVideoInfo(videoId)
	if err == nil {
		v.Data["json"] = ReturnSuccess(200, "success", video, 0)
	} else {
		v.Data["json"] = ReturnError(4004, "没有相关内容")
	}
	v.ServeJSON()
}

// @Summer 获取剧集
func (v *AdvertController) VideoEpisodesList() {
	videoId, _ := v.GetInt("videoId")
	if videoId == 0 {
		v.Data["json"] = ReturnError(4001, "必须指定视频ID")
	}
	num, episodes, err := models.GetVideoEpisodeList(videoId)
	if err == nil {
		v.Data["json"] = ReturnSuccess(200, "success", episodes, num)
	} else {
		v.Data["json"] = ReturnError(4004, "没有相关内容")
	}
	v.ServeJSON()
}
