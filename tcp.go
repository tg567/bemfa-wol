package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

func tcpWOL() {
	defer func() {
		if r := recover(); r != nil {
			println(r)
		}
	}()
	var con net.Conn
	var err error
	var connectTime int32
	for {
		con, err = net.DialTimeout("tcp", "bemfa.com:8344", 5*time.Second)
		if err != nil {
			println("tcp连接错误", err)
			time.Sleep(time.Minute * 5 * time.Duration(connectTime))
			if connectTime < 6 {
				connectTime++
			}
			continue
		}
		defer con.Close()

		connectTime = 0

		// 处理连接
		ctx, cancel := context.WithCancel(context.Background())
		go heartbeat(ctx, con)
		if _, err := con.Write([]byte(fmt.Sprintf("cmd=1&uid=%s&topic=%s\r\n", uid, topic))); err != nil {
			println("订阅topic错误", err)
			cancel()
			return
		}
		handleConnection(con)
		cancel()
	}
}

// 处理单个连接的函数
func handleConnection(con net.Conn) {

	// 这里可以添加具体的逻辑来处理连接
	reader := bufio.NewReader(con)

	for {
		lineBytes, _, err := reader.ReadLine()
		if err != nil {
			println("tcp read错误", err)
			return
		}
		line := string(lineBytes)
		line = strings.ReplaceAll(line, "\r", "")
		line = strings.ReplaceAll(line, "\n", "")
		values, err := url.ParseQuery(line)
		if err != nil {
			println("解析参数错误", err)
			continue
		}
		switch values.Get("cmd") {
		case "2":
			if values.Get("topic") == topic {
				switch values.Get("msg") {
				case "on":
					wol()
				case "off":
					output, err := exec.Command("ssh", sshUserServer, `shutdown`, `-s`, `-t`, `0`).Output()
					if err != nil {
						println("ssh shutdown错误", err)
					}
					if string(output) != "" {
						println("ssh shutdown output:", string(output))
					}
				}
			}
			println("返回参数", line)
		}
	}

}

// 处理单个连接的函数
func heartbeat(ctx context.Context, con net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			println(r)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			println("done heartbeat")
			return
		default:
			time.Sleep(time.Minute)
			_, err := con.Write([]byte("ping\r\n"))
			if err != nil {
				println("heartbeat错误", err)
			}
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

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:9", broadcastAddress))
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
