package main

import (
	"flag"
	"fmt"
)
var (
	logPath string
	bindPort int
)

func init() {
	flag.StringVar(&logPath, "l", "", "设置日志路径")
	flag.IntVar(&bindPort, "p", 8889, "设置绑定端口")
	flag.Parse()
}
func main() {
	fmt.Println(logPath)
	fmt.Println(bindPort)
}