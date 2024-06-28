package main

import (
	"os/exec"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func mqttWOL() {
	clientID := uid
	broker := "bemfa.com:9501"
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetConnectTimeout(time.Second * 3)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetMaxReconnectInterval(5 * time.Second)
	opts.SetClientID(clientID)
	opts.SetKeepAlive(time.Minute)
	opts.SetProtocolVersion(4)

	client := mqtt.NewClient(opts)

	token := client.Connect()
	if !token.WaitTimeout(time.Second * 3) {
		println("mqtt连接 3秒超时")
		return
	}
	if err := token.Error(); err != nil {
		println("token错误", err)
		return
	}

	client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		defer func() {
			if r := recover(); r != nil {
				println(r)
			}
		}()
		println("收到消息", string(msg.Payload()))
		switch string(msg.Payload()) {
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
	})

}
