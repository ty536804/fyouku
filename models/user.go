package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type User struct {
	Id       int
	Name     string
	Mobile   string
	Password string
	Status   int
	AddTime  int64
	Avatar   string
}

func init() {
	orm.RegisterModel(new(User))
}

// @Summer 判断手机号是否已经注册
func IsUserMobile(mobile string) bool {
	o := orm.NewOrm()
	user := User{Mobile: mobile}
	err := o.Read(&user, "Mobile")
	if err == orm.ErrNoRows {
		return false
	} else if err == orm.ErrMissPK {
		return false
	}
	return true
}

//@ Summer保存用户
func UserSave(mobile, password string) error {
	o := orm.NewOrm()
	user := User{
		Name:     "",
		Mobile:   mobile,
		Password: password,
		Status:   1,
		AddTime:  time.Now().Unix(),
		Avatar:   "",
	}
	_, err := o.Insert(&user)
	return err
}

func IsMobileLogin(mobile, password string) (int, string) {
	o := orm.NewOrm()
	user := User{}
	err := o.QueryTable("user").Filter("mobile", mobile).
		Filter("password", password).One(&user)
	if err == orm.ErrNoRows {
		return 0, ""
	} else if err == orm.ErrMissPK {
		return 0, ""
	}
	return user.Id, user.Name
}
