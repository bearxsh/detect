package main

import (
	"fmt"
	"runtime"
)

func main()  {
	goos := runtime.GOOS
	fmt.Println(goos)
	goarch := runtime.GOARCH
	fmt.Println(goarch)

}
