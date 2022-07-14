package main

import (
	"fmt"
	"time"
)

func main()  {
	var tempDelay time.Duration // how long to sleep on accept failure
	for  {
		if tempDelay == 0 {
			tempDelay = 5 * time.Millisecond
		} else {
			tempDelay *= 2
		}
		if max := 1 * time.Second; tempDelay > max {
			tempDelay = max
		}
		fmt.Printf("http: Accept error; retrying in %v\n", tempDelay)
		time.Sleep(tempDelay)
	}

}
