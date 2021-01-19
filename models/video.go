package models

import "github.com/astaxie/beego/orm"

type VideoData struct {
	Id            int
	Title         string
	SubTitle      string
	AddTime       int64
	Img           string
	Img1          string
	EpisodesCount int
	IsEnd         int
}

type Video struct {
	Id int
	Title string
	SubTitle string
	AddTime int64
	Img string
	Img1 string
	EpisodesCount int
	IsEnd int
	ChannelId int
	Status int
	RegionId int
	TypeId int
	Sort int
	EpisodesUpdateTime int
	Comment int
}

func init()  {
	orm.RegisterModel(new(Video))
}

func GetChannelHotList(channelId int) (int64, []VideoData, error) {
	o := orm.NewOrm()

	var videos []VideoData
	num, err := o.Raw("SELECT title,sub_title,add_time,img,img1,episodes_count,is_end FROM video WHERE status=1 AND is_hot=1 AND channel_id=?" +
		" ORDER BY episodes_update_time DESC LIMIT 9",channelId).QueryRows(&videos)
	return num,videos,err
}

func GetChannelRecommendRegionList(channelId,regionId int) (int64, []VideoData, error)  {
	o := orm.NewOrm()
	var videos []VideoData
	num, err := o.Raw("SELECT title,sub_title,add_time,img,img1,episodes_count,is_end FROM video WHERE status=1 AND is_recommend=1 AND region_id = ? AND channel_id=?" +
		" ORDER BY episodes_update_time DESC LIMIT 9",regionId,channelId).QueryRows(&videos)
	return num,videos,err
}

func GetChannelRecommendTypeList(channelId,typeId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var videos []VideoData
	num, err := o.Raw("SELECT title,sub_title,add_time,img,img1,episodes_count,is_end FROM video WHERE status=1 AND is_recommend=1 AND type_id = ? AND channel_id=?" +
		" ORDER BY episodes_update_time DESC LIMIT 9",typeId,channelId).QueryRows(&videos)
	return num,videos,err
}

func GetChannelVideoList(channelId,regionId,typeId,limit,offset int,end,sort string)  (int64, []orm.Params, error) {
	o := orm.NewOrm()
	var videos []orm.Params

	qs := o.QueryTable("video")
	qs = qs.Filter("channel_id",channelId)
	qs = qs.Filter("status",1)
	if regionId > 0 {
		qs = qs.Filter("region_id",regionId)
	}
	if typeId > 0 {
		qs = qs.Filter("type_id", typeId)
	}
	if end == "n" {//未完结
		qs = qs.Filter("is_emd",0)
	} else  if end == "y" {//已完结
		qs = qs.Filter("is_emd",1)
	}
	//剧集更新时间排序
	if sort == "episodesUpdate" { //倒序
		qs = qs.OrderBy("-episodes_update_time")
	} else if sort == "comment" {
		qs = qs.OrderBy("-comment")
	} else if sort == "addTime" {
		qs = qs.OrderBy("-add_time")
	} else  {
		qs = qs.OrderBy("-add_time")
	}
	nums, _ := qs.Values(&videos,"id","title","sub_title","add_time","img","img1","episodes_count","is_end")
	qs = qs.Limit(limit,offset)
	_,err := qs.Values(&videos,"id","title","sub_title","add_time","img","img1","episodes_count","is_end")
	return nums,videos,err
}