package main

import (
	"fmt"
	"github.com/prometheus/procfs"
)

func main()  {
	fs, err := procfs.NewFS("/proc")
	if err != nil {
		panic(err)
	}
	fmt.Println(fs.CPUInfo())
}
