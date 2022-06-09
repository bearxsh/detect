package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"time"
)

func init()  {
	// 初始化日志
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
}
func main() {
	ticker := time.NewTicker(time.Duration(3) * time.Second)
	defer ticker.Stop()
	// 判断不通的ip会阻塞很久，这里使用-t参数设置超时时间
	param := "https://39.96.21"
	command := fmt.Sprintf("httping -c 5 -t 1 %s | awk 'END{print $4}' | awk -F/ '{print $2}'", param)
	for {
		select {
		case <-ticker.C:
		}
		// 注意：Output末尾会有换行符 \n
		output, err := exec.Command("/bin/sh", "-c", command).Output()
		if err != nil {
			log.Errorf("Failed to exec command [%s]: %s", command, err)
			// 只能通过控制器改变任务状态
			//taskChannelMap.Delete(task.TaskId)
			//db.Model(&DetectTask{}).Where("task_id=?", task.TaskId).Update("status", 0)
			//return
		}
		result := string(output)
		log.Infof(" result [%s]", result[:len(result)-1])
	}
}
