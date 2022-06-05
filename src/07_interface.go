package main

import (
	"fmt"
	"net"
)

func main()  {
	a, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(a); i++ {
		fmt.Println(a[i].Name, a[i].Flags)

	}

}
