package main

import (
	"fmt"
	"fyoukuApi/Services/mq"
)

func main() {
	mq.Consumer("", "fyouku_demo", callback)
}

func callback(s string) {
	fmt.Printf("msg is:%s \n", s)
}
