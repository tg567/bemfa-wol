package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
	"time"
)

func tcpWOL() {
	con, err := net.DialTimeout("tcp", "bemfa.com:8344", 5*time.Second)
	if err != nil {
		log.Println("tcp连接错误", err)
		return
	}
	defer con.Close()

	ch := make(chan struct{})
	// 处理连接
	go handleConnection(ch, con)
	go heartbeat(con)
	if _, err := con.Write([]byte(fmt.Sprintf("cmd=1&uid=%s&topic=%s\r\n", uid, topic))); err != nil {
		log.Println("订阅topic错误", err)
		return
	}

	<-ch
	time.Sleep(time.Second)
}

// 处理单个连接的函数
func handleConnection(ch chan struct{}, con net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	// 这里可以添加具体的逻辑来处理连接
	reader := bufio.NewReader(con)

	for {
		select {
		case <-ch:
			return
		default:
			lineBytes, _, err := reader.ReadLine()
			if err != nil {
				close(ch)
				log.Println("tcp read错误", err)
				return
			}
			line := string(lineBytes)
			line = strings.ReplaceAll(line, "\r", "")
			line = strings.ReplaceAll(line, "\n", "")
			values, err := url.ParseQuery(line)
			if err != nil {
				log.Println("解析参数错误", err)
				continue
			}
			switch values.Get("cmd") {
			case "2":
				if values.Get("topic") == topic && values.Get("msg") == "on" {
					wol()
				}
				fallthrough
			default:
				log.Println("返回参数", line)
			}
		}

	}
}

// 处理单个连接的函数
func heartbeat(con net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	for {
		log.Println("别关，小爱同学唤醒电脑用的！！！！！！")
		time.Sleep(time.Minute)
		_, err := con.Write([]byte("ping\r\n"))
		if err != nil {
			log.Println("heartbeat错误", err)
		}
	}
}

func wol() {
	byteArray := make([]byte, 102)
	for i := 0; i < 6; i++ {
		byteArray[i] = 0xFF
	}

	mac = strings.ReplaceAll(mac, ":", "")
	mac = strings.ReplaceAll(mac, "-", "")
	macBytes, err := hex.DecodeString(mac)
	if err != nil {
		panic("mac地址错误")
	}

	for i := 1; i <= 16; i++ {
		for j := 0; j < 6; j++ {
			byteArray[i*6+j] = macBytes[j]
		}
	}

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:9", ipAddress))
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	_, err = conn.Write(byteArray)
	if err != nil {
		panic(err)
	}
}
