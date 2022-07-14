package main

import (
	"fmt"
	"time"
)

func main() {

	ticker := time.NewTicker(10 * time.Second) // 设置轮询时间
	defer ticker.Stop()
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"))
	for {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"))
		<-ticker.C

	}
}
