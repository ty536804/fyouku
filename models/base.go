package models

import "github.com/astaxie/beego/orm"

type Region struct {
	Id int
	Name string
}

func GetChannelRegion(channelId int) (int64, []Region, error)  {
	o := orm.NewOrm()
	var regions []Region

	nums, err := o.Raw("SELECT id,name FROM channel_region WHERE status=1 AND channel_id= ? ORDER BY sort DESC",channelId).QueryRows(&regions)
	return nums,regions,err
}

type Type struct {
	Id int
	Name string
}

func GetChannelType(channelId int) (int64, []Type, error)  {
	o := orm.NewOrm()
	var regions []Type

	nums, err := o.Raw("SELECT id,name FROM channel_type WHERE status=1 AND channel_id= ? ORDER BY sort DESC",channelId).QueryRows(&regions)
	return nums,regions,err
}