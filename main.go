package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

// uid
var uid = "xxx"

// topic
var topic = "xxx"

// mac地址
var mac = "00:00:00:00:00:00"

// 机器网段
var ipAddress = "192.168.1.255"

var logFile string
var file *os.File

func main() {
	flag.StringVar(&logFile, "f", "", "log file path")
	flag.Parse()
	if logFile != "" {
		var err error
		file, err = os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			log.Println("日志文件路径不存在")
			return
		}
	}

	//tcp网络唤醒
	tcpWOL()
	//mqtt网络唤醒
	// mqttWOL()
}

func println(a ...any) {
	log.Println(a...)
	if file != nil {
		b := []any{time.Now().Format("2006-01-02 15:04:05")}
		b = append(b, a...)
		fmt.Fprintln(file, b...)
	}
}
