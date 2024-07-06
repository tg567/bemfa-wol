package main

import (
	"github.com/tg567/bemfa-wol/gui"
	"github.com/tg567/bemfa-wol/utils"
)

func main() {
	utils.LoadConfig()
	w := gui.LoadWindow()
	w.ShowAndRun()
	// flag.StringVar(&uid, "uid", "", "bemfa uid")
	// flag.StringVar(&topic, "topic", "", "topic")
	// flag.StringVar(&mac, "mac", "", "mac")
	// flag.StringVar(&broadcastAddress, "broadcast", "", "broadcast address")
	// flag.StringVar(&sshUserServer, "ssh", "", "ssh user@ipaddress")
	// flag.StringVar(&logFile, "f", "", "log file path")
	// flag.StringVar(&wolType, "type", "tcp", "wol type, tcp/mqtt")
	// flag.Parse()
	// if logFile != "" {
	// 	var err error
	// 	file, err = os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	// 	if err != nil {
	// 		utils.Println("日志文件路径不存在")
	// 		return
	// 	}
	// }

	// if uid == "" || topic == "" || mac == "" || sshUserServer == "" || broadcastAddress == "" {
	// 	utils.Println("参数错误")
	// 	return
	// }

	// if wolType == "tcp" {
	// 	//tcp网络唤醒
	// 	go tcpWOL()
	// 	utils.Println("tcpWOL start...")
	// } else {
	// 	//mqtt网络唤醒
	// 	mqttWOL()
	// 	utils.Println("mqttWOL start...")
	// }

	// ch := make(chan struct{})
	// <-ch
}
