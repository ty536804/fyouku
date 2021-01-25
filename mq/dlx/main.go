package main

import (
	"fmt"
	"fyoukuApi/Services/mq"
)

func main() {
	mq.ConsumerDlx("fyouku.dlx.a", "fyouku_dlx_a", "fyouku.dlx.b", "fyouku_dlx.b", 10000, callback)
}

func callback(s string) {
	fmt.Printf("msg is:%s \n", s)
}
