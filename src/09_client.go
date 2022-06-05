package main

import (
	"encoding/binary"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func init() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		// 注意：2006-01-02 15:04:05是固定的，不能改动！
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			// 精简文件名
			fileName := path.Base(frame.File)
			return frame.Function, fileName + ":" + strconv.Itoa(frame.Line)
		},
	})
}

func main() {
	conn, err := net.DialTimeout("tcp", "localhost:20880", 3 * time.Second)
	defer conn.Close()
	if err != nil {
		log.Errorf("connect error: %s", err)
		return
	}
	log.Info("连接成功")
	//content := []byte("Hello worldHello world我是，佛挡杀佛")
/*	data := []byte{0x00, 0x00, 0x00, 0x12, 0x02, 0x00, 0x00, 0x00, 0x00,0x00,0x00,0x00,0x03, 'h','e','l','l','o',0x00}
	err = binary.Write(conn, binary.BigEndian, data)
	data = []byte{0x00, 0x00, 0x12, 0x02, 0x00, 0x00, 0x00, 0x00,0x00,0x00,0x00,0x03, 'h','e','l','l','d', 0x00}
	err = binary.Write(conn, binary.BigEndian, data)
	time.Sleep(2 * time.Second)
	data = []byte{0x00, 0x00, 0x12, 0x02, 0x00, 0x00, 0x00, 0x00,0x00,0x00,0x00,0x03, 'w','o','r','l','d'}
	err = binary.Write(conn, binary.BigEndian, data)*/

	data := []byte{0xda, 0xbb, 0xc0}
	err = binary.Write(conn, binary.BigEndian, data)
/*	data = []byte{0x12, 0x02, 0x00, 0x00, 0x00, 0x00,0x00,0x00,0x00,0x03, 'h','e','l','l','o'}
	err = binary.Write(conn, binary.BigEndian, data)*/
	//n, err := conn.Write(data)
	if err != nil {
		log.Errorf("write error: %s", err)
	}
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	fmt.Printf("读取到的字节数：%d", n)
	fmt.Printf("返回结果：%s", string(buf[:n]))
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
	//log.Infof("write %d bytes", n)

}
