package main

import (
	"fmt"
	"time"
)

func main() {
	//var c chan int
	c := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		c <- 3
		time.Sleep(3 * time.Second)
		c <- 4
		time.Sleep(3 * time.Second)
		c <- 5
	}()

	for item := range c {
		fmt.Println(item)
	}
}
