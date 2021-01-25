package main

import (
	"fmt"
	"fyoukuApi/Services/mq"
)

func main() {
	//获取所有
	mq.ConsumeEx("fyouku.demo.topic", "topic", "fyouku.*", callback)
}

func callback(s string) {
	fmt.Printf("msg is:%s \n", s)
}
