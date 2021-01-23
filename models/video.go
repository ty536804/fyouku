package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

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
	Id                 int
	Title              string
	SubTitle           string
	AddTime            int64
	Img                string
	Img1               string
	EpisodesCount      int
	IsEnd              int
	ChannelId          int
	Status             int
	RegionId           int
	TypeId             int
	Sort               int
	EpisodesUpdateTime int64
	Comment            int
	UserId             int
}

type Episodes struct {
	Id      int
	Title   string
	AddTime int64
	Num     int
	PlayUrl string
	Comment int
}

func init() {
	orm.RegisterModel(new(Video))
}

func GetChannelHotList(channelId int) (int64, []VideoData, error) {
	o := orm.NewOrm()

	var videos []VideoData
	num, err := o.Raw("SELECT title,sub_title,add_time,img,img1,episodes_count,is_end FROM video WHERE status=1 AND is_hot=1 AND channel_id=?"+
		" ORDER BY episodes_update_time DESC LIMIT 9", channelId).QueryRows(&videos)
	return num, videos, err
}

func GetChannelRecommendRegionList(channelId, regionId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var videos []VideoData
	num, err := o.Raw("SELECT title,sub_title,add_time,img,img1,episodes_count,is_end FROM video WHERE status=1 AND is_recommend=1 AND region_id = ? AND channel_id=?"+
		" ORDER BY episodes_update_time DESC LIMIT 9", regionId, channelId).QueryRows(&videos)
	return num, videos, err
}

func GetChannelRecommendTypeList(channelId, typeId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var videos []VideoData
	num, err := o.Raw("SELECT title,sub_title,add_time,img,img1,episodes_count,is_end FROM video WHERE status=1 AND is_recommend=1 AND type_id = ? AND channel_id=?"+
		" ORDER BY episodes_update_time DESC LIMIT 9", typeId, channelId).QueryRows(&videos)
	return num, videos, err
}

func GetChannelVideoList(channelId, regionId, typeId, limit, offset int, end, sort string) (int64, []orm.Params, error) {
	o := orm.NewOrm()
	var videos []orm.Params

	qs := o.QueryTable("video")
	qs = qs.Filter("channel_id", channelId)
	qs = qs.Filter("status", 1)
	if regionId > 0 {
		qs = qs.Filter("region_id", regionId)
	}
	if typeId > 0 {
		qs = qs.Filter("type_id", typeId)
	}
	if end == "n" { //未完结
		qs = qs.Filter("is_emd", 0)
	} else if end == "y" { //已完结
		qs = qs.Filter("is_emd", 1)
	}
	//剧集更新时间排序
	if sort == "episodesUpdate" { //倒序
		qs = qs.OrderBy("-episodes_update_time")
	} else if sort == "comment" {
		qs = qs.OrderBy("-comment")
	} else if sort == "addTime" {
		qs = qs.OrderBy("-add_time")
	} else {
		qs = qs.OrderBy("-add_time")
	}
	nums, _ := qs.Values(&videos, "id", "title", "sub_title", "add_time", "img", "img1", "episodes_count", "is_end")
	qs = qs.Limit(limit, offset)
	_, err := qs.Values(&videos, "id", "title", "sub_title", "add_time", "img", "img1", "episodes_count", "is_end")
	return nums, videos, err
}

func GetVideoInfo(videoId int) (Video, error) {
	o := orm.NewOrm()
	var video Video
	err := o.Raw("SELECT * FROM video WHERE id = ? Limit 1", videoId).QueryRow(&video)
	return video, err
}

func GetVideoEpisodeList(videoId int) (int64, []Video, error) {
	o := orm.NewOrm()
	var video []Video
	nums, err := o.Raw("SELECT id,title,add_time,num,play_url,comment FROM video_episodes WHERE video_id= ? AND status=1 ORDER BY num ASC", videoId).QueryRows(&video)
	return nums, video, err
}

func GetChannelTop(channelId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var videos []VideoData
	num, err := o.Raw("SELECT id,title,sub_title,img,img1,add_time,episodes_count,is_end FROM video WHERE channel_id=? AND status=1 ORDER BY comment DESC LIMIT 10", channelId).QueryRows(&videos)
	return num, videos, err
}

func GetTypeTop(typeId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var videos []VideoData
	num, err := o.Raw("SELECT id,title,sub_title,img,img1,add_time,episodes_count,is_end FROM video WHERE type_id=? AND status=1 ORDER BY comment DESC LIMIT 10", typeId).QueryRows(&videos)
	return num, videos, err
}

func GetUserVideo(uid int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var videos []VideoData
	nums, err := o.Raw("id,title,sub_title,img,img1,add_time,episodes_count,is_end FROM video WHERE user_id=? ORDER BY add_time DESC", uid).QueryRows(&videos)
	return nums, videos, err
}

func SaveVideo(playUrl, title, subTitle, aliyunVideoId string, channelId, typeId, regionId, uid int) error {
	times := time.Now().Unix()
	o := orm.NewOrm()
	var video Video
	video.Title = title
	video.SubTitle = subTitle
	video.AddTime = times
	video.Img = ""
	video.Img1 = ""
	video.Status = 1
	video.ChannelId = channelId
	video.TypeId = typeId
	video.RegionId = regionId
	video.Comment = 0
	video.EpisodesUpdateTime = times
	video.UserId = uid
	videoId, err := o.Insert(&video)
	if err == nil {
		if aliyunVideoId != "" {
			playUrl = ""
		}
		o.Raw("INSERT INTO video_episodes (title,add_time,num,video_id,play_url,"+
			"status,comment,aliyun_video_id) VALUES (?,?,?,?,?,?,?)", subTitle, times, 1, videoId, playUrl, 1, 0, aliyunVideoId).Exec()
	}
	return err
}

func SaveAliYunVideo(videoId, log string) error {
	o := orm.NewOrm()
	_, err := o.Raw("INSERT INTO aliyun_video (video_id,log,add_time) VALUES(?,?,?)", videoId, log, time.Now().Unix()).Exec()
	return err
}
