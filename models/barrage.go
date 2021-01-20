package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Barrage struct {
	Id          int
	Content     string
	CurrentTime int
	AddTime     int64
	UserId      int
	Status      int
	EpisodesId  int
	VideoId     int
}

type BarrageData struct {
	Id          int    `json:"id"`
	Content     string `json:"content"`
	CurrentTime int    `json:"current_time"`
}

func init() {
	orm.RegisterModel(new(Barrage))
}

func BarrageList(episodesId, startTime, endTime int) (int64, []BarrageData, error) {
	o := orm.NewOrm()
	var barrage []BarrageData
	nums, err := o.Raw("SELECT id,content,current_time FROM barrage WHERE status=1 AND episodes_id =? AND `current_time>=?` AND `current_time<?` ORDER BY `current_time` ASC", episodesId, startTime, endTime).QueryRows(&barrage)
	return nums, barrage, err
}

func SaveBarrage(episodesId, videoId, currentTime, userId int, content string) error {
	o := orm.NewOrm()
	var barrage Barrage
	barrage.CurrentTime = currentTime
	barrage.EpisodesId = episodesId
	barrage.UserId = userId
	barrage.Status = 1
	barrage.VideoId = videoId
	barrage.Content = content
	barrage.AddTime = time.Now().Unix()
	_, err := o.Insert(&barrage)
	return err
}
