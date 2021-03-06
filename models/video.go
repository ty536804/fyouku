package models

import (
	"encoding/json"
	"fyoukuApi/Services/Es"
	redisClient "fyoukuApi/Services/redis"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"strconv"
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
	Comment       int
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
	IsRecommend        int
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

//根据传入参数获取视频列表
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

//根据传入参数获取视频列表
func GetChannelVideoListEs(channelId, regionId, typeId, limit, offset int, end, sort string) (int64, []Video, error) {
	query := make(map[string]interface{})
	bools := make(map[string]interface{})
	var must []map[string]interface{}
	must = append(must, map[string]interface{}{"term": map[string]interface{}{"channel_id": channelId}})
	must = append(must, map[string]interface{}{"term": map[string]interface{}{"status": 1}})
	if regionId > 0 {
		must = append(must, map[string]interface{}{"term": map[string]interface{}{"region_id": regionId}})
	}
	if typeId > 0 {
		must = append(must, map[string]interface{}{"term": map[string]interface{}{"type_id": typeId}})
	}
	if end == "n" { //未完结
		must = append(must, map[string]interface{}{"term": map[string]interface{}{"is_emd": 0}})
	} else if end == "y" { //已完结
		must = append(must, map[string]interface{}{"term": map[string]interface{}{"is_emd": 1}})
	}
	bools["must"] = must
	query["bool"] = bools
	sortData := []map[string]string{map[string]string{"add_time": "desc"}}
	//剧集更新时间排序
	if sort == "episodesUpdate" { //倒序
		sortData = []map[string]string{map[string]string{"episodes_update_time": "desc"}}
	} else if sort == "comment" {
		sortData = []map[string]string{map[string]string{"comment": "desc"}}
	} else if sort == "addTime" {
		sortData = []map[string]string{map[string]string{"add_time": "desc"}}
	}
	res := Es.EsSearch("fyouku_video", query, offset, limit, sortData)
	total := res.Total.Value
	var data []Video
	for _, v := range res.Hits {
		var itemData Video
		err := json.Unmarshal([]byte(v.Source), &itemData)
		if err != nil {
			data = append(data, itemData)
		}
	}
	return int64(total), data, nil
}

func GetVideoInfo(videoId int) (Video, error) {
	o := orm.NewOrm()
	var video Video
	err := o.Raw("SELECT * FROM video WHERE id = ? Limit 1", videoId).QueryRow(&video)
	return video, err
}

//增加redis缓存 获取视频详情
func RedisGetVideoInfo(videoId int) (Video, error) {
	var video Video
	conn := redisClient.PoolConnect()
	defer conn.Close()
	//定义redis key
	redisKey := "video:id:" + string(videoId)
	// 判断key中是否存在
	exists, err := redis.Bool(conn.Do("exists", redisKey))
	if exists {
		res, _ := redis.Values(conn.Do("hgetall", redisKey))
		err = redis.ScanStruct(res, &video)
	} else {
		o := orm.NewOrm()
		err := o.Raw("SELECT * FROM video WHERE id = ? Limit 1", videoId).QueryRow(&video)
		if err == nil {
			_, err := conn.Do("hmset", redis.Args{redisKey}.AddFlat(video)...)
			if err == nil {
				conn.Do("expire", redisKey, 86400)
			}
		}
	}
	return video, err
}

// 获取视频剧集列表
func GetVideoEpisodeList(videoId int) (int64, []Video, error) {
	o := orm.NewOrm()
	var video []Video
	nums, err := o.Raw("SELECT id,title,add_time,num,play_url,comment FROM video_episodes WHERE video_id= ? AND status=1 ORDER BY num ASC", videoId).QueryRows(&video)
	return nums, video, err
}

// @获取视频剧集列表
func RedisGetVideoEpisodeList(videoId int) (int64, []Episodes, error) {
	var (
		episodes []Episodes
		num      int64
		err      error
	)

	conn := redisClient.PoolConnect()
	defer conn.Close()

	redisKey := "video:episode:videoId:" + strconv.Itoa(videoId)
	exists, err := redis.Bool(conn.Do("exists", redisKey))
	if exists {
		num, err = redis.Int64(conn.Do("llen", redisKey))
		if err == nil {
			values, _ := redis.Values(conn.Do("lrange", redisKey, "0", "-1"))
			var episodesInfo Episodes
			for _, v := range values {
				err = json.Unmarshal(v.([]byte), &episodesInfo)
				if err == nil {
					episodes = append(episodes, episodesInfo)
				}
			}
		}
	} else {
		o := orm.NewOrm()
		num, err = o.Raw("SELECT id,title,add_time,num,play_url,comment FROM video_episodes WHERE video_id= ? AND status=1 ORDER BY num ASC", videoId).QueryRows(&episodes)
		if err == nil {
			//变量获取到的信息，把信息json化保存
			for v, _ := range episodes {
				jsonVal, err := json.Marshal(v)
				if err == nil {
					conn.Do("rpush", redis.Args{redisKey}.AddFlat(jsonVal))
				}
			}
			conn.Do("expire", redisKey, 86400)
		}
	}
	return num, episodes, err
}

// 频道排行榜
func GetChannelTop(channelId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var videos []VideoData
	num, err := o.Raw("SELECT id,title,sub_title,img,img1,add_time,episodes_count,is_end FROM video WHERE channel_id=? AND status=1 ORDER BY comment DESC LIMIT 10", channelId).QueryRows(&videos)
	return num, videos, err
}

// redis 频道排行榜
func RedisGetChannelTop(channelId int) (int64, []VideoData, error) {
	var (
		videos []VideoData
		err    error
		num    int64
	)

	conn := redisClient.PoolConnect()
	defer conn.Close()

	redisKey := "video:top:channel:channelId" + strconv.Itoa(channelId)

	exists, err := redis.Bool(conn.Do("exists", redisKey))
	if exists {
		num = 0
		res, _ := redis.Values(conn.Do("zrevrange", redisKey, "0", "10", "WITHSCORES"))
		for k, v := range res {
			//1属性  2 分数
			if k%2 == 0 {
				videoId, err := strconv.Atoi(string(v.([]byte)))
				videoInfo, err := RedisGetVideoInfo(videoId)
				if err == nil {
					var videoDataInfo VideoData
					videoDataInfo.Id = videoInfo.Id
					videoDataInfo.Img = videoInfo.Img
					videoDataInfo.Img1 = videoInfo.Img1
					videoDataInfo.IsEnd = videoInfo.IsEnd
					videoDataInfo.SubTitle = videoInfo.SubTitle
					videoDataInfo.Title = videoInfo.Title
					videoDataInfo.AddTime = videoInfo.AddTime
					videoDataInfo.Comment = videoInfo.Comment
					videoDataInfo.EpisodesCount = videoInfo.EpisodesCount
					videos = append(videos, videoDataInfo)
					num++
				}
			}
		}
	} else {
		o := orm.NewOrm()
		num, err = o.Raw("SELECT id,title,sub_title,img,img1,add_time,episodes_count,is_end,comment FROM video WHERE channel_id=? AND status=1 ORDER BY comment DESC LIMIT 10", channelId).QueryRows(&videos)
		if err == nil {
			//保存redis
			for _, v := range videos {
				conn.Do("zadd", redisKey, v.Comment, v.Id)
			}
			conn.Do("expire", redisKey, 86400*30)
		}
	}
	return num, videos, err
}

// @Summer 获取类型排行榜
func GetTypeTop(typeId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var videos []VideoData
	num, err := o.Raw("SELECT id,title,sub_title,img,img1,add_time,episodes_count,is_end FROM video WHERE type_id=? AND status=1 ORDER BY comment DESC LIMIT 10", typeId).QueryRows(&videos)
	return num, videos, err
}

//redis 频道类型排行榜
func RedisGetTypeTop(typeId int) (int64, []VideoData, error) {
	var (
		videos []VideoData
		err    error
		num    int64
	)

	conn := redisClient.PoolConnect()
	defer conn.Close()

	redisKey := "redis:top:type:typeId" + strconv.Itoa(typeId)
	exists, err := redis.Bool(conn.Do("exists", redisKey))
	if exists {
		num = 0
		res, _ := redis.Values(conn.Do("zrevrange", redisKey, "0", "10", "WITHSCORES"))
		for k, v := range res {
			//1属性  2 分数
			if k%2 == 0 {
				videoId, err := strconv.Atoi(string(v.([]byte)))
				videoInfo, err := RedisGetVideoInfo(videoId)
				if err == nil {
					var videoDataInfo VideoData
					videoDataInfo.Id = videoInfo.Id
					videoDataInfo.Img = videoInfo.Img
					videoDataInfo.Img1 = videoInfo.Img1
					videoDataInfo.IsEnd = videoInfo.IsEnd
					videoDataInfo.SubTitle = videoInfo.SubTitle
					videoDataInfo.Title = videoInfo.Title
					videoDataInfo.AddTime = videoInfo.AddTime
					videoDataInfo.Comment = videoInfo.Comment
					videoDataInfo.EpisodesCount = videoInfo.EpisodesCount
					videos = append(videos, videoDataInfo)
					num++
				}
			}
		}
	} else {
		o := orm.NewOrm()
		num, err = o.Raw("SELECT id,title,sub_title,img,img1,add_time,episodes_count,is_end FROM video WHERE type_id=? AND status=1 ORDER BY comment DESC LIMIT 10", typeId).QueryRows(&videos)
		if err == nil {
			//保存redis
			for _, v := range videos {
				conn.Do("zadd", redisKey, v.Comment, v.Id)
			}
			conn.Do("expire", redisKey, 86400*30)
		}
	}
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

//获取所有视频
func GetAllList() (int64, []Video, error) {
	o := orm.NewOrm()
	var videos []Video
	num, err := o.Raw("SELECT id,title,sub_title,status,add_time,img,img1,channel_id,type_id," +
		"region_id,user_id,episodes_count,episodes_update_time,is_end,is_hot,is_recommend,comment FROM video").QueryRows(&videos)
	return num, videos, err
}
