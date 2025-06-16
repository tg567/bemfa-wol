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

func tcpWOL(devices []Device, uid string) {
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

		connectTime = 0

		// 处理连接
		ctx, cancel := context.WithCancel(context.Background())
		go heartbeat(ctx, con)
		deviceMap := make(map[string]Device, len(devices))
		for _, v := range devices {
			deviceMap[v.Topic] = v
			if _, err := con.Write([]byte(fmt.Sprintf("cmd=1&uid=%s&topic=%s\r\n", uid, v.Topic))); err != nil {
				println("订阅topic错误", err)
				con.Close()
				cancel()
				return
			}
		}

		reader := bufio.NewReader(con)
		for {
			lineBytes, _, err := reader.ReadLine()
			if err != nil {
				println("tcp read错误", err)
				break
			}
			line := string(lineBytes)
			line = strings.ReplaceAll(line, "\r", "")
			line = strings.ReplaceAll(line, "\n", "")
			values, err := url.ParseQuery(line)
			if err != nil {
				println("解析参数错误", err)
				continue
			}
			device, ok := deviceMap[values.Get("topic")]
			if values.Get("cmd") == "2" && ok {
				handleDeviceOperation(&device, values.Get(`msg`))
			}
		}
		con.Close()
		cancel()
	}
}

func handleDeviceOperation(device *Device, msg string) {
	println(msg, device)
	switch msg {
	case "on":
		wol(device)
	case "off":
		output, err := exec.Command("ssh", fmt.Sprintf(`%s@%s`, device.User, device.IP), `shutdown`, `-s`, `-t`, `0`).Output()
		if err != nil {
			println("ssh shutdown错误", err)
		}
		if string(output) != "" {
			println("ssh shutdown output:", string(output))
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
	t := time.NewTicker(time.Minute)
	for {
		select {
		case <-ctx.Done():
			println("done heartbeat")
			return
		case <-t.C:
			_, err := con.Write([]byte("ping\r\n"))
			if err != nil {
				println("heartbeat错误", err)
			}
		}
	}
}

func wol(device *Device) {
	byteArray := make([]byte, 102)
	for i := 0; i < 6; i++ {
		byteArray[i] = 0xFF
	}

	mac := device.MAC
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

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:9", device.Broadcast))
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
