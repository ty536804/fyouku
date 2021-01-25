package models

import (
	"encoding/json"
	"fyoukuApi/Services/mq"
	"github.com/astaxie/beego/orm"
	"time"
)

type Comment struct {
	Id          int
	Content     string
	AddTime     int64
	UserId      int
	Stamp       int
	Status      int
	PraiseCount int
	EpisodesId  int
	VideoId     int
}

func init() {
	orm.RegisterModel(new(Comment))
}

func GetCommentList(episodesId, offset, limit int) (int64, []Comment, error) {
	o := orm.NewOrm()
	var comments []Comment
	nums, err := o.Raw("SELECT id FROM comment WHERE status=1 AND episodes_id=?", episodesId).QueryRows(&comments)
	_, err = o.Raw("SELECT id,content,add_time,user_id,praise_count,episodes_id FROM comment WHERE status=1 AND episodes_id=? ORDER BY add_time DESC LIMIT ? OFFSET ?", episodesId, limit, offset).QueryRows(&comments)
	return nums, comments, err
}

func SaveComment(content string, uid, episodesId, videoId int) error {
	o := orm.NewOrm()
	var comment Comment
	comment.Content = content
	comment.UserId = uid
	comment.EpisodesId = episodesId
	comment.VideoId = videoId
	comment.Status = 1
	comment.Stamp = 0
	comment.AddTime = time.Now().Unix()
	_, err := o.Insert(&comment)
	if err == nil {
		//修改视频的总评论数
		o.Raw("UPDATE video SET comment=comment+1 WHERE id=?", videoId).Exec()
		//修改视频剧集的评论数
		o.Raw("UPDATE video_episodes SET comment=comment+1 WHERE id=?", episodesId).Exec()
		//创建一个简单的模式的MQ 把要传递的数据转换为json字符串
		videoObj := map[string]int{"VideoId": videoId}
		videoJson, _ := json.Marshal(videoObj)
		mq.Publish("", "fyouku_top", string(videoJson))
	}
	return err
}
