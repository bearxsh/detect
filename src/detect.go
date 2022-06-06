package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type DetectTask struct {
	Id       string
	Action   string
	Type     string
	Param    string
	Interval int
	Status   int
}

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

var db *sql.DB
var insertStmt *sql.Stmt
var updateStmt *sql.Stmt
var taskChannelMap sync.Map

func init() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		// 注意：2006-01-02 15:04:05.000 是固定的，不能改动！
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			// 精简文件名
			fileName := path.Base(frame.File)
			return frame.Function, fileName + ":" + strconv.Itoa(frame.Line)
		},
	})

	var err error
	db, err = sql.Open("sqlite3", "./detector.db")
	if err != nil {
		log.Errorf("connect database fail: %s", err)
		panic(err)
	}
	log.Info("connect database success")

	sqlStmt := `CREATE TABLE IF NOT EXISTS "detect_task" (
		"id" INTEGER NOT NULL,
		"action" TEXT NOT NULL,
		"type" TEXT NOT NULL,
		"param" TEXT NOT NULL,
		"interval" TEXT NOT NULL,
		"status" INTEGER NOT NULL,
		PRIMARY KEY ("id")
	);`
	if _, err := db.Exec(sqlStmt); err != nil {
		log.Errorf("sqlStmt execute fail: %s", err)
		panic(err)
	}

	insertStmt, err = db.Prepare("INSERT INTO detect_task(id, action, type, param, interval, status) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Errorf("generate insert prepare fail: %s", err)
		panic(err)
	}

	updateStmt, err = db.Prepare("UPDATE detect_task SET status=? WHERE id=?")
	if err != nil {
		log.Errorf("generate update prepare fail: %s", err)
		panic(err)
	}

}

// 在defer中设置数据库任务状态为停止
func ping(quit chan struct{}, addr string, interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	// 判断不通的ip会阻塞 使用-t参数设置超时时间
	//Command := fmt.Sprintf("ping -c 1 -t 1 101.42.110.140 > /dev/null && echo true || echo false")
	Command := fmt.Sprintf("ping -c 1 -t 1 %s > /dev/null && echo true || echo false", addr)
	for {
		select {
		case <-quit:
			log.Infof("task quit")
			return
		case <-ticker.C:

		}
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
		/*	_, err = conn.Write(resString)
			if err != nil {
				log.Errorf("write err: %s", err)
				return
			}*/
		address := "localhost:9999"
		// open connection to server
		client, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err != nil {
			log.Errorf("connect err: %s", err)
			// 处理连接失败情况
			continue
		}
		_, err = client.Write(resString)
		if err != nil {
			log.Errorf("write err: %s", err)
			//continue
		}
		err = client.Close()
		if err != nil {
			log.Errorf("close err: %s", err)
			//return
		}
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

	// 短连接其实可以不用for循环的
	for {
		// 确保buf容量足够
		buf := make([]byte, 4*1024)
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
				// 保存检测任务
				_, err := insertStmt.Exec(detectTask.Id, detectTask.Action, detectTask.Type, detectTask.Param, detectTask.Interval, 1)
				if err != nil {

					log.Errorf("insert err: %s", err)

					res := Response{301, err.Error()}
					result, _ := json.Marshal(res)
					// 返回处理结果
					_, err = conn.Write(result)
					if err != nil {
						log.Errorf("write err: %s", err)
						return
					}
					return
				}

				quit := make(chan struct{})
				taskChannelMap.Store(detectTask.Id, quit)
				go ping(quit, detectTask.Param, detectTask.Interval)
			}

		} else if detectTask.Action == "stop" {
			_, err := updateStmt.Exec(0, detectTask.Id)
			if err != nil {
				log.Errorf("update err: %s", err)
				res := Response{301, err.Error()}
				result, _ := json.Marshal(res)
				// 返回处理结果
				_, err = conn.Write(result)
				if err != nil {
					log.Errorf("write err: %s", err)
					return
				}
				return
			}
			value, exist := taskChannelMap.LoadAndDelete(detectTask.Id)
			if exist {
				ch, ok := value.(chan struct{})
				if ok {
					ch <- struct{}{}
				} else {
					log.Warnf("convert fail")
				}

			} else {
				log.Warnf("stop task, but task [%s] is not exist or not running.", detectTask.Id)
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

//  重启时恢复正在执行的任务
func loadTaskOnStart()  {

	rows, err := db.Query("SELECT * FROM detect_task WHERE status=1")
	if err != nil {
		log.Errorf("select err: %s", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var task DetectTask
		err = rows.Scan(&task.Id, &task.Action, &task.Type, &task.Param, &task.Interval, &task.Status)
		if err != nil {
			log.Errorf("iteration err: %s", err)
		}
		quit := make(chan struct{})
		taskChannelMap.Store(task.Id, quit)
		go ping(quit, task.Param, task.Interval)
	}
	err = rows.Err()
	if err != nil {
		log.Errorf("iteration err: %s", err)
	}
}
// 命令行参数指定运行端口
func main() {

	loadTaskOnStart()

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
