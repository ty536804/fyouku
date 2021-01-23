package controllers

import (
	"encoding/json"
	"fmt"
	"fyoukuApi/models"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vod"
	"github.com/astaxie/beego"
)

type AliYunController struct {
	beego.Controller
}

var (
	accessKeyId     = "LTAI4GKK942rberZCkydxdtN"
	accessKeySecret = "Bb2jB4NEhM5835ynHuTRGwpRSeq6ao"
)

type JSONS struct {
	RequestId     string
	UploadAddress string
	UploadAuth    string
	VideoId       string
}

func (a *AliYunController) InitVodClient(accessKeyId string, accessKeySecret string) (client *vod.Client, err error) {

	// 点播服务接入区域
	regionId := "cn-shanghai"

	// 创建授权对象
	credential := &credentials.AccessKeyCredential{
		accessKeyId,
		accessKeySecret,
	}

	// 自定义config
	config := sdk.NewConfig()
	config.AutoRetry = true     // 失败是否自动重试
	config.MaxRetryTime = 3     // 最大重试次数
	config.Timeout = 3000000000 // 连接超时，单位：纳秒；默认为3秒

	// 创建vodClient实例
	return vod.NewClientWithOptions(regionId, config, credential)
}

func (a *AliYunController) MyCreateUploadVideo(client *vod.Client, title, desc, filename, coverUrl, tags string) (response *vod.CreateUploadVideoResponse, err error) {
	request := vod.CreateCreateUploadVideoRequest()

	request.Title = title
	request.Description = desc
	request.FileName = filename
	//request.CateId = "-1"
	request.CoverURL = coverUrl
	request.Tags = tags

	request.AcceptFormat = "JSON"
	return client.CreateUploadVideo(request)
}

// @Summer aliyun/create/upload/video
func (a *AliYunController) CreateUploadVideo() {
	title := a.GetString("title")
	desc := a.GetString("desc")
	filename := a.GetString("filename")
	coverUrl := a.GetString("coverUrl")
	tags := a.GetString("tags")
	client, err := a.InitVodClient(accessKeyId, accessKeySecret)
	if err != nil {
		panic(err)
	}

	response, err := a.MyCreateUploadVideo(client, title, desc, filename, coverUrl, tags)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.GetHttpContentString())
	data := &JSONS{
		response.VideoId,
		response.RequestId,
		response.UploadAddress,
		response.UploadAuth,
	}
	a.Data["json"] = data
	a.ServeJSON()
	//fmt.Printf("VideoId: %s\n UploadAddress: %s\n UploadAuth: %s",
	//	response.VideoId, response.UploadAddress, response.UploadAuth)
}

func (a *AliYunController) MyRefreshUploadVideo(client *vod.Client, VideoId string) (response *vod.RefreshUploadVideoResponse, err error) {
	request := vod.CreateRefreshUploadVideoRequest()
	request.VideoId = VideoId
	request.AcceptFormat = "JSON"

	return client.RefreshUploadVideo(request)
}

// @Summer aliyun/refresh/upload/video
func (a *AliYunController) main() {
	videoId := a.GetString("videoId")
	client, err := a.InitVodClient(accessKeyId, accessKeySecret)
	if err != nil {
		panic(err)
	}

	response, err := a.MyRefreshUploadVideo(client, videoId)
	if err != nil {
		panic(err)
	}
	data := &JSONS{
		response.VideoId,
		response.RequestId,
		response.UploadAddress,
		response.UploadAuth,
	}
	a.Data["json"] = data
	a.ServeJSON()
	//fmt.Println(response.GetHttpContentString())
	//fmt.Printf("UploadAddress: %s\n UploadAuth: %s", response.UploadAddress, response.UploadAuth)
}

func (a *AliYunController) MyGetPlayAuth(client *vod.Client, videoId string) (response *vod.GetVideoPlayAuthResponse, err error) {
	request := vod.CreateGetVideoPlayAuthRequest()
	request.VideoId = videoId
	request.AcceptFormat = "JSON"

	return client.GetVideoPlayAuth(request)
}

type PlayJSONS struct {
	PlayAuth string
}

// @Summer aliyun/video/play/auth
func (a *AliYunController) GetPlayAuth() {
	videoId := a.GetString("videoId")
	client, err := a.InitVodClient(accessKeyId, accessKeySecret)
	if err != nil {
		panic(err)
	}

	response, err := a.MyGetPlayAuth(client, videoId)
	if err != nil {
		panic(err)
	}

	data := &PlayJSONS{
		response.PlayAuth,
	}
	a.Data["json"] = data
	a.ServeJSON()
	//fmt.Println(response.GetHttpContentString())
	//fmt.Printf("%s: %s\n", response.VideoMeta, response.PlayAuth)
}

type CallBackData struct {
	EventTime   string
	EventType   string
	VideoId     string
	Status      int
	Exteng      string
	StreamInfos []CallBackStreamInfosData
}

type CallBackStreamInfosData struct {
	Status     string
	Bitrate    int
	Definition string
	Duration   int
	Encrypt    bool
	FileUrl    string
	Format     string
	Fps        int
	Height     int
	Size       int
	Width      int
	JobId      int
}

//回调函数路由
// aliyun/video/callback
func (a *AliYunController) VideoCallBack() {
	var ob CallBackData
	r := a.Ctx.Input.RequestBody //获取Body信息  json格式
	json.Unmarshal(r, &ob)
	models.SaveAliYunVideo(ob.VideoId, string(r))
	a.Ctx.WriteString("success")
}
