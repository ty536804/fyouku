package main

import (
	"fmt"
	"fyoukuApi/Services/mq"
)

func main() {
	mq.ConsumeEx("fyouku.demo.direct", "direct", "two", callback)
}

func callback(s string) {
	fmt.Printf("msg is:%s \n", s)
}
