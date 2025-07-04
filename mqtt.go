package main

import (
	"context"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func mqttWOL(devices []Device, uid string) {
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

	ipMap := make(map[string]Device)

	for _, v := range devices {
		if v.IP != "" {
			ipMap[v.IP] = v
		}
	}

	pHandler := newMQTTPingHandler(uid, ipMap, client)
	go pHandler.ping(context.Background(), false)

	for i := range devices {
		client.Subscribe(devices[i].Topic, 1, func(client mqtt.Client, msg mqtt.Message) {
			defer func() {
				if r := recover(); r != nil {
					println(r)
				}
			}()
			println("收到消息", string(msg.Payload()))
			handleDeviceOperation(&devices[i], string(msg.Payload()))
		})
	}
}
