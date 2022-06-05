package main

// import necessary packages
import (
	"fmt"
	"net"
	"time"
)

func main() {

	// initialize the address of the server
	address := "101.42.110.40:22"

	// open connection to server
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)

	// check if connection was successfully established
	if err != nil {
		fmt.Println("The following error occured", err)
	} else {
		fmt.Println("The connection was established to", conn)
	}
	//conn.Close()

}
