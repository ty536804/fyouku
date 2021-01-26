package models

import (
	redisClient "fyoukuApi/Services/redis"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
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

type UserInfo struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	AddTime int64  `json:"add_time"`
	Avatar  string `json:"avatar"`
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

//根据用户ID获取用户信息
func GetUserInfo(uid int) (UserInfo, error) {
	o := orm.NewOrm()
	var user UserInfo
	err := o.Raw("SELECT id,name,add_time,avatar FROM user WHERE id=? LIMIT 1", uid).QueryRow(&user)
	return user, err
}

// @Title 获取用户信息
func RedisGetUserInfo(uid int) (UserInfo, error) {
	var user UserInfo

	conn := redisClient.PoolConnect()
	defer conn.Close()
	redisKey := "user:uid:" + string(uid)
	exists, err := redis.Bool(conn.Do("exists", redisKey))
	if exists {
		res, _ := redis.Values(conn.Do("hgetall", redisKey))
		err = redis.ScanStruct(res, &user)
	} else {
		o := orm.NewOrm()
		err := o.Raw("SELECT id,name,add_time,avatar FROM user WHERE id=? LIMIT 1", uid).QueryRow(&user)
		if err == nil {
			_, err := conn.Do("hmset", redis.Args{redisKey}.AddFlat(user))
			if err == nil {
				conn.Do("expire", redisKey, 86400)
			}
		}
	}
	return user, err
}
