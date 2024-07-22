package connection

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/tg567/bemfa-wol/utils"
)

var ErrorClose = fmt.Errorf("close")

func TcpWOL(stopChan <-chan struct{}) {
	defer func() {
		if r := recover(); r != nil {
			utils.Println(r)
		}
	}()
	utils.Println("tcp wol starting...")
	var con net.Conn
	var err error
	var connectTime int32
	for {
		con, err = net.DialTimeout("tcp", "bemfa.com:8344", 5*time.Second)
		if err != nil {
			utils.Println("tcp连接错误", err)
			time.Sleep(time.Minute * 5 * time.Duration(connectTime))
			if connectTime < 6 {
				connectTime++
			}
			continue
		}

		connectTime = 0

		// 处理连接
		go heartbeat(stopChan, con)
		if _, err := con.Write([]byte(fmt.Sprintf("cmd=1&uid=%s&topic=%s\r\n", utils.WolConfig.UID, utils.WolConfig.Topic))); err != nil {
			con.Close()
			utils.Println("订阅topic错误", err)
			return
		}
		err := handleConnection(stopChan, con)
		if err == ErrorClose {
			con.Close()
			utils.Println("done TcpWOL")
			return
		}
		con.Close()
	}
}

// 处理单个连接的函数
func handleConnection(stopChan <-chan struct{}, con net.Conn) error {
	// 这里可以添加具体的逻辑来处理连接
	reader := bufio.NewReader(con)

	for {
		select {
		case <-stopChan:
			utils.Println("done handleConnection")
			return ErrorClose
		default:
			lineBytes, _, err := reader.ReadLine()
			if err != nil {
				utils.Println("tcp read错误", err)
				return err
			}
			line := string(lineBytes)
			line = strings.ReplaceAll(line, "\r", "")
			line = strings.ReplaceAll(line, "\n", "")
			values, err := url.ParseQuery(line)
			if err != nil {
				utils.Println("解析参数错误", err)
				continue
			}
			switch values.Get("cmd") {
			case "2":
				if values.Get("topic") == utils.WolConfig.Topic {
					switch values.Get("msg") {
					case "on":
						wol()
					case "off":
						output, err := exec.Command("ssh", utils.WolConfig.SSH, `shutdown`, `-s`, `-t`, `0`).Output()
						if err != nil {
							utils.Println("ssh shutdown错误", err)
						}
						if string(output) != "" {
							utils.Println("ssh shutdown output:", string(output))
						}
					}
				}
				utils.Println("返回参数", line)
			}
		}
	}
}

// 处理单个连接的函数
func heartbeat(stopChan <-chan struct{}, con net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			utils.Println(r)
		}
	}()
	for {
		select {
		case <-stopChan:
			utils.Println("done heartbeat")
			return
		default:
			time.Sleep(time.Minute)
			_, err := con.Write([]byte("ping\r\n"))
			if err != nil {
				utils.Println("heartbeat错误", err)
			}
		}
	}
}

func wol() {
	byteArray := make([]byte, 102)
	for i := 0; i < 6; i++ {
		byteArray[i] = 0xFF
	}

	mac := strings.ReplaceAll(utils.WolConfig.Mac, ":", "")
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

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:9", utils.WolConfig.Broadcast))
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
