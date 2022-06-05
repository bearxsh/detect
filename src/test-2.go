package main

import (
	"fmt"
	"net"
	"os"
)
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
func main()  {
	fmt.Fprintf(os.Stderr, "result:\n")
	address, err := net.ResolveTCPAddr("tcp4", "101.42.110.40:220")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, address)
	checkError(err)
	err = conn.Close()
	checkError(err)



}
