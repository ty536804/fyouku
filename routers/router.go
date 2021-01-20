// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"fyoukuApi/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/register/save", &controllers.UserController{}, "get:SaveRegister")
	beego.Router("/login/do", &controllers.UserController{}, "get:LoginDo")
	beego.Router("/channel/advert", &controllers.AdvertController{}, "get:ChannelAdvert")
	beego.Router("/chanel/hot/list", &controllers.AdvertController{}, "get:ChanelHotList")
	beego.Router("/channel/recommend/region/list", &controllers.AdvertController{}, "get:ChannelRecommendRegionList")
	beego.Router("/channel/recommend/type/list", &controllers.AdvertController{}, "get:GetChannelRecommendTypeList")
	beego.Router("/channel/region", &controllers.BaseController{}, "get:ChannelRegion")
	beego.Router("/channel/type", &controllers.BaseController{}, "get:ChannelType")
	beego.Router("/channel/video", &controllers.BaseController{}, "get:ChannelVideo")
	beego.Router("/video/info", &controllers.AdvertController{}, "get:VideoInfo")
	beego.Router("/video/episodes/list", &controllers.AdvertController{}, "get:VideoEpisodesList")
	beego.Router("/comment/list", &controllers.CommentController{}, "get:List")
	beego.Router("/comment/save", &controllers.CommentController{}, "get:Save")
	beego.Router("/channel/top", &controllers.TopController{}, "get:ChannelTop")
	beego.Router("/type/top", &controllers.TopController{}, "get:TypeTop")
	beego.Router("/send/message/do", &controllers.UserController{}, "get:SendMessageDo")
	//beego.Include(&controllers.UserController{})
	//ns := beego.NewNamespace("/v1",
	//	beego.NSNamespace("/object",
	//		beego.NSInclude(
	//			&controllers.ObjectController{},
	//		),
	//	),
	//	beego.NSNamespace("/user",
	//		beego.NSInclude(
	//			&controllers.UserController{},
	//		),
	//	),
	//)
	//beego.AddNamespace(ns)
}
