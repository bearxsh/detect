package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main()  {
	 name := os.Args[1]
	 fmt.Println(name)
	//name := "eth0"
	command := fmt.Sprintf("ethtool %s | awk 'END{print $3}'", name)
	// 注意：Output末尾会有换行符 \n
	output, err := exec.Command("/bin/sh", "-c", command).Output()
	if err != nil {
		fmt.Printf("Failed to exec command [%s]: %s\n", command, err)
		return
	}
	result := string(output)
    fmt.Println(len(result))
	res := result[:len(result)-1]
	fmt.Printf("result [%s]\n", res)
	if res != "yes" && res != "no" {
		fmt.Println("unknown result")
	}
	if res == "yes" {
		fmt.Println("on")
	} else if res == "no" {
		fmt.Println("off")
	}
}