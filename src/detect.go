package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"time"
)

type DetectTask struct {
	Id string
	Action string
	Type string
	Param string
	Interval int
}

type Response struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
}

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

func ping(conn net.Conn, addr string, interval int)  {
	// 判断不通的ip会阻塞 使用-t参数设置超时时间
	//Command := fmt.Sprintf("ping -c 1 -t 1 101.42.110.140 > /dev/null && echo true || echo false")
	Command := fmt.Sprintf("ping -c 1 -t 1 %s > /dev/null && echo true || echo false", addr)
	for {
		// 注意：Output末尾会有换行符 \n
		output, err := exec.Command("/bin/sh", "-c", Command).Output()
		if err != nil {
			log.Errorf("shell exec err: %s", err)
			return
		}
		result := string(output)
		/*	parseBool, err := strconv.ParseBool(result[:len(result)-1])
			if err != nil {
				log.Errorf("string parse to bool err: %s", err)
				return
			}*/
		log.Infof("ping [%s], result [%s]", addr, result[:len(result)-1])
		res := Response{200, result[:len(result)-1]}
		resString, _ := json.Marshal(res)
		// 返回处理结果
		_, err = conn.Write(resString)
		if err != nil {
			log.Errorf("write err: %s", err)
			return
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}



}

func handleConn(conn net.Conn) {

	defer func(conn net.Conn) {
		log.Infof("close conn [%s]", conn.RemoteAddr())
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
		log.Infof("read %d bytes, content is %s", n, string(buf[:n]))
		detectTask := new(DetectTask)
		err = json.Unmarshal(buf[:n], detectTask)
		if err != nil {
			log.Errorf("json decode err: %s, input: %s", err, string(buf[:n]))
			return
		}
		log.Info("json decode success: ", detectTask)


		// TODO 处理业务逻辑
		if detectTask.Action == "create" {
			if detectTask.Type == "ping" {
				go ping(conn, detectTask.Param, detectTask.Interval)
			}

		}


		res := Response{200, "ok"}
		result, _ := json.Marshal(res)
		// 返回处理结果
		_, err = conn.Write(result)
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
