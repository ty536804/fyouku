package controllers

import (
	"fmt"
	redisClient "fyoukuApi/Services/redis"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
)

type RedisDemoController struct {
	beego.Controller
}

func (r *RedisDemoController) Demo() {
	c := redisClient.PoolConnect()
	defer c.Close()

	_, err := c.Do("SET", "username", "frog")
	if err == nil {
		//设置过期时间
		c.Do("expire", "username", 1000)
	}

	val, err := redis.String(c.Do("get", "username"))
	if err == nil {
		fmt.Println(1)
		fmt.Println(val)
		//获取剩余过期时间  获取到的值是int64
		ttl, _ := redis.Int64(c.Do("ttl", "username"))
		fmt.Println(ttl)
	} else {
		fmt.Println(22)
		fmt.Println(err)
	}
	r.Ctx.WriteString(val)
}
