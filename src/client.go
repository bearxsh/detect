package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
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
func isCanceled(ch chan bool) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}
func pingIP(ch chan bool) {
	log.Trace("Something very low level.")
	log.Debug("Useful debugging information.")
	log.Info("Something noteworthy happened!")
	log.Warn("You should probably take a look at this.")
	log.Error("Something failed but I'm not quitting.")
	count := 0
	for {
		Command := fmt.Sprintf("ping -c 1 -t 1 101.42.110.140 > /dev/null && echo true || echo false")
		// 注意：Output末尾会有换行符\n
		output, _ := exec.Command("/bin/sh", "-c", Command).Output()
		result := string(output)
		xx, _ := strconv.ParseBool(result[:len(result)-1])
		fmt.Println(xx)
		if xx {
			count = 0
		} else {
			count++
		}
		if count == 3 {
			fmt.Println("exit")
			break
		}
		if isCanceled(ch) {
			fmt.Println("任务取消")
			break
		} else {
			time.Sleep(3 * time.Second)
		}
		//bools, success := <-ch
		//fmt.Print(err)

	}
}

func main() {
	ch := make(chan bool, 1)
	go pingIP(ch)
	time.Sleep(10 * time.Second)
	close(ch)
	// 通过WaitGroup防止主线程退出 这会导致goroutine也退出
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
