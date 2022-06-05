package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// 判断不通的ip会阻塞 使用-t参数
	//Command := fmt.Sprintf("ping -c 1 -t 1 101.42.110.140 > /dev/null && echo true || echo false")
	Command := fmt.Sprintf("ping -c 1 -t 1 39.96.210.14 > /dev/null && echo true || echo false")
	output, err := exec.Command("/bin/sh", "-c", Command).Output()
	fmt.Print(string(output))
	fmt.Print(err)
}