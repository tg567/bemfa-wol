package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var file *os.File
var config *Setting

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "", "config file path")
	flag.Parse()

	var err error
	config, err = loadSetting(configPath)
	if err != nil {
		println("打开配置文件错误", err)
	}

	if config.LogFile != "" {
		var err error
		file, err = os.OpenFile(config.LogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			println("日志文件路径不存在")
			return
		}
	}

	if err := config.Validate(); err != nil {
		println("参数错误", err)
		return
	}

	if strings.ToLower(config.Type) == "tcp" {
		//tcp网络唤醒
		go tcpWOL(config.Devices, config.UID)
		println("tcpWOL start...")
	} else {
		//mqtt网络唤醒
		mqttWOL(config.Devices, config.UID)
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
