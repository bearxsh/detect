package main

import (
	"fmt"
	"time"
)

func main() {

	ticker := time.NewTicker(3 * time.Second) // 设置轮询时间
	defer ticker.Stop()

	for {

		<-ticker.C
		fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"))
	}
}
