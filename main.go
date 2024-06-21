package main

// uid
var uid = "xxx"

// topic
var topic = "xxx"

// mac地址
var mac = "00:00:00:00:00:00"

// 机器网段
var ipAddress = "192.168.1.255"

func main() {
	//tcp网络唤醒
	tcpWOL()
	//mqtt网络唤醒
	// mqttWOL()
}
