package controllers

import (
	"encoding/json"
	"fyoukuApi/Services/Es"
	"fyoukuApi/models"
	"github.com/astaxie/beego"
	"strconv"
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

//搜索接口
func (v *AdvertController) Search() {
	keyWord := v.GetString("keyword")
	//获取翻页信息
	limit, _ := v.GetInt("limit")
	offset, _ := v.GetInt("offse")
	if keyWord == "" {
		v.Data["json"] = ReturnError(4001, "关键字不能为空")
		v.ServeJSON()
	}
	if limit == 0 {
		limit = 12
	}
	sott := []map[string]string{map[string]string{
		"id": "desc",
	}}
	query := map[string]interface{}{
		"bool": map[string]interface{}{
			"must": map[string]interface{}{
				"term": map[string]interface{}{
					"title": keyWord,
				},
			},
		},
	}
	res := Es.EsSearch("fyouku_video", query, offset, limit, sott)
	total := res.Total.Value
	var data []models.Video
	for _, v := range res.Hits {
		var itemData models.Video
		err := json.Unmarshal([]byte(v.Source), &itemData)
		if err != nil {
			data = append(data, itemData)
		}
	}
	if total > 0 {
		v.Data["json"] = ReturnSuccess(200, "success", data, int64(total))
	} else {
		v.Data["json"] = ReturnError(4004, "没有相关内容")
	}
	v.ServeJSON()
}

//导入es脚本
func (v *AdvertController) SendEs() {
	_, data, _ := models.GetAllList()
	for _, v := range data {
		body := map[string]interface{}{
			"id":                   v.Id,
			"title":                v.Title,
			"sub_title":            v.SubTitle,
			"add_time":             v.AddTime,
			"img":                  v.Img,
			"img1":                 v.Img1,
			"episodes_count":       v.EpisodesCount,
			"is_end":               v.IsEnd,
			"channel_id":           v.ChannelId,
			"status":               v.Status,
			"region_id":            v.RegionId,
			"type_id":              v.TypeId,
			"episodes_update_time": v.EpisodesUpdateTime,
			"comment":              v.Comment,
			"user_id":              v.UserId,
			"is_recommend":         v.IsRecommend,
		}
		Es.EsAdd("fyouku_video", "video-"+strconv.Itoa(v.Id), body)
	}
}
