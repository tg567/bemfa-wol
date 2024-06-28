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
var broadcastAddress = "192.168.1.255"

// ssh用户服务器
var sshUserServer = "root@192.168.1.1"

var logFile string
var file *os.File
var wolType string

func main() {
	flag.StringVar(&uid, "uid", "", "bemfa uid")
	flag.StringVar(&topic, "topic", "", "topic")
	flag.StringVar(&mac, "mac", "", "mac")
	flag.StringVar(&broadcastAddress, "broadcast", "", "broadcast address")
	flag.StringVar(&sshUserServer, "ssh", "", "ssh user@ipaddress")
	flag.StringVar(&logFile, "f", "", "log file path")
	flag.StringVar(&wolType, "type", "tcp", "wol type, tcp/mqtt")
	flag.Parse()
	if logFile != "" {
		var err error
		file, err = os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			println("日志文件路径不存在")
			return
		}
	}

	if uid == "" || topic == "" || mac == "" || sshUserServer == "" || broadcastAddress == "" {
		println("参数错误")
		return
	}

	if wolType == "tcp" {
		//tcp网络唤醒
		go tcpWOL()
		println("tcpWOL start...")
	} else {
		//mqtt网络唤醒
		mqttWOL()
		println("mqttWOL start...")
	}

	ch := make(chan struct{})
	<-ch
}

func println(a ...any) {
	log.Println(a...)
	if file != nil {
		b := []any{time.Now().Format("2006-01-02 15:04:05")}
		b = append(b, a...)
		fmt.Fprintln(file, b...)
	}
}
