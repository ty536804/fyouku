package main

import (
	"fmt"
	"fyoukuApi/Services/mq"
)

func main() {
	mq.ConsumeEx("fyouku.demo.fanout", "fanout", "", callback)
}

func callback(s string) {
	fmt.Printf("msg is:%s \n", s)
}
