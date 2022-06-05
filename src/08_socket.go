package main

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"path"
	"runtime"
	"strconv"
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

func handleConn(conn net.Conn) {

	defer func(conn net.Conn) {
		log.Info("close conn...")
		err := conn.Close()
		if err != nil {
			log.Errorf("conn close, error: %s", err)
		}
	}(conn)

	for {
		buf := make([]byte, 4 * 1024)
		err := conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if err != nil {
			log.Errorf("conn SetReadDeadline error: %s", err)
			return
		}
		n, err := conn.Read(buf)
		if err == io.EOF {
			log.Warn("eof occur")
			break
		}
		if err != nil {
			log.Errorf("conn read %d bytes, error: %s", n, err)
			//if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			//	// 进行其他业务逻辑的处理
			//	continue
			//}
			return
		}
		log.Infof("read %d bytes, content is %s\n", n, string(buf[:n]))
		// TODO 处理业务逻辑
		// 发送
		_, err = conn.Write(buf[:n])
		if err != nil {
			log.Errorf("write err: %s", err)
			return
		}

	}
}
func main() {
	listen, err := net.Listen("tcp", ":8889")
	if err != nil {
		log.Errorf("listen error: %s", err)
		return
	}
	log.Info("server starts success")
	for {
		accept, err := listen.Accept()
		if err != nil {
			log.Errorf("accept error: %s", err)
			break
		}
		log.Infof("remote connect success: %s", accept.RemoteAddr())
		go handleConn(accept)
	}
}
