package main

import (
	"fmt"
	"os"
)

func main() {
	fileName := "/etc/logrotate.d/logrotate-test"
	f,err := os.Create(fileName)
	defer f.Close()
	if err !=nil {
		fmt.Println(err.Error())
	} else {
		_,err=f.Write([]byte("要写入的文本内容"))
		if err != nil {
			panic(err)
		}
	}
}
