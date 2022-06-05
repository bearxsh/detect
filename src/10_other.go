package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Runner interface {
	Run()
}

type Program struct {
	/* fields */
	name string
}

func (p Program) Run() {
	fmt.Print("fd", p.name)
}

func newProgrm(p *Program) {
	p.name = " wang"
}

func main() {
	/*	fileName:="/Users/ruby/Documents/pro/a/english.txt"
		file,err := os.Open(fileName)
		if err != nil{
			fmt.Println(err)
			return
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		var buf []byte
		reader.Read(buf)
		//binary.Write()
		var p Program
		p.name = " zhang"
		newProgrm(new(Program))
		p.Run()*/
	data := []byte{0x00, 0x00, 0x00}
	bytes.NewReader(data)
	var totalLen int32
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &totalLen)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(totalLen)
	


}
