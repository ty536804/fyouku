package controllers

import (
	"fyoukuApi/models"
	"github.com/astaxie/beego"
)

type CommentController struct {
	beego.Controller
}

type CommentInfo struct {
	Id           int             `json:"id"`
	Content      string          `json:"content"`
	AddTime      int64           `json:"add_time"`
	AddTimeTitle string          `json:"add_time_title"`
	UserId       int             `json:"user_id"`
	Stamp        int             `json:"stamp"`
	PraiseCount  int             `json:"praise_count"`
	UserInfo     models.UserInfo `json:"user_info"`
}

// 获取评论列表
//func (c *CommentController) List() {
//	//获取剧集数
//	episodesId, _ := c.GetInt("episodesId")
//	// 获取页码信息
//	limit, _ := c.GetInt("limit")
//	offset, _ := c.GetInt("offset")
//	if episodesId == 0 {
//		c.Data["json"] = ReturnError(4001, "必须指定视频剧集")
//		c.ServeJSON()
//	}
//	if limit == 0 {
//		limit = 12
//	}
//	num, comments, err := models.GetCommentList(episodesId, offset, limit)
//	if err == nil {
//		var data []CommentInfo
//		var commentInfo CommentInfo
//		for _, v := range comments {
//			commentInfo.Id = v.Id
//			commentInfo.AddTime = v.AddTime
//			commentInfo.AddTimeTitle = DateFormat(v.AddTime)
//			commentInfo.Content = v.Content
//			commentInfo.UserId = v.UserId
//			commentInfo.Stamp = v.Stamp
//			commentInfo.PraiseCount = v.PraiseCount
//			user, _ := models.RedisGetUserInfo(v.UserId)
//			commentInfo.UserInfo = user
//			data = append(data, commentInfo)
//		}
//		c.Data["json"] = ReturnSuccess(200, "success", data, num)
//	} else {
//		c.Data["json"] = ReturnError(4004, "没有相关内容")
//	}
//	c.ServeJSON()
//}

// 获取评论列表
func (c *CommentController) List() {
	//获取剧集数
	episodesId, _ := c.GetInt("episodesId")
	// 获取页码信息
	limit, _ := c.GetInt("limit")
	offset, _ := c.GetInt("offset")
	if episodesId == 0 {
		c.Data["json"] = ReturnError(4001, "必须指定视频剧集")
		c.ServeJSON()
	}
	if limit == 0 {
		limit = 12
	}
	num, comments, err := models.GetCommentList(episodesId, offset, limit)
	if err == nil {
		var data []CommentInfo
		var commentInfo CommentInfo
		// 获取uid channel
		uidChan := make(chan int, 12)
		resChan := make(chan models.UserInfo, 12)
		closeChan := make(chan bool, 5)
		// 把获取到的uid放到channel中
		go func() {
			for _, v := range comments {
				uidChan <- v.UserId
			}
			close(uidChan)
		}()
		//处理uidChan中的信息
		for i := 0; i < 5; i++ {
			go chanGetUserInfo(uidChan, resChan, closeChan)
		}
		//判断是否执行完成，信息集合
		go func() {
			for i := 0; i < 5; i++ {
				<-closeChan
			}
			close(resChan)
			close(closeChan)
		}()

		userInfoMap := make(map[int]models.UserInfo)
		for r := range resChan {
			userInfoMap[r.Id] = r
		}

		for _, v := range comments {
			commentInfo.Id = v.Id
			commentInfo.AddTime = v.AddTime
			commentInfo.AddTimeTitle = DateFormat(v.AddTime)
			commentInfo.Content = v.Content
			commentInfo.UserId = v.UserId
			commentInfo.Stamp = v.Stamp
			commentInfo.PraiseCount = v.PraiseCount
			commentInfo.UserInfo = userInfoMap[v.UserId]
			data = append(data, commentInfo)
		}

		c.Data["json"] = ReturnSuccess(200, "success", data, num)
	} else {
		c.Data["json"] = ReturnError(4004, "没有相关内容")
	}
	c.ServeJSON()
}

func chanGetUserInfo(uidChan chan int, resChan chan models.UserInfo, closeChan chan bool) {
	for uid := range uidChan {
		res, err := models.GetUserInfo(uid)
		if err == nil {
			resChan <- res
		}
	}
	closeChan <- true
}

func (c *CommentController) Save() {
	content := c.GetString("content")
	uid, _ := c.GetInt("uid")
	episodesId, _ := c.GetInt("episodesId")
	videoId, _ := c.GetInt("videoId")
	if content == "" {
		c.Data["json"] = ReturnError(4001, "内容不能为空")
		c.ServeJSON()
	}
	if uid == 0 {
		c.Data["json"] = ReturnError(4002, "请先登录")
		c.ServeJSON()
	}
	if episodesId == 0 {
		c.Data["json"] = ReturnError(4003, "必须指定剧集ID")
		c.ServeJSON()
	}
	if videoId == 0 {
		c.Data["json"] = ReturnError(4005, "必须指定视频ID")
		c.ServeJSON()
	}
	err := models.SaveComment(content, uid, episodesId, videoId)
	if err == nil {
		c.Data["json"] = ReturnSuccess(200, "success", 0, 1)
	} else {
		c.Data["json"] = ReturnError(5000, err)
	}
	c.ServeJSON()
}
