package main

// import necessary packages
import (
	"fmt"
	"net"
)

func main() {

	// initialize the address of the server
	address := "39.96.210.14:6688"

	// open connection to server
	//conn, err := net.DialTimeout("udp", address, 2*time.Second)
	conn, err := net.Dial("udp", address)
	if err != nil {
		fmt.Println("连接UDP服务器失败，err: ", err)
		return
	}
	write, err := conn.Write([]byte("fdfsd"))
	fmt.Printf("write: %d\n", write)
	// check if connection was successfully established
	if err != nil {
		fmt.Println("The following error occured", err)
	} else {
		fmt.Println("The connection was established to", conn)
		conn.Close()
	}

}
