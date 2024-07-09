package connection

import (
	"os/exec"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tg567/bemfa-wol/utils"
)

func MqttWOL(closeChan <-chan struct{}) {
	defer func() {
		if r := recover(); r != nil {
			utils.Println(r)
		}
	}()
	utils.Println("mqtt wol starting...")
	clientID := utils.WolConfig.UID
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
		utils.Println("mqtt连接 3秒超时")
		return
	}
	if err := token.Error(); err != nil {
		utils.Println("token错误", err)
		return
	}
	defer token.Done()

	client.Subscribe(utils.WolConfig.Topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		defer func() {
			if r := recover(); r != nil {
				utils.Println(r)
			}
		}()
		utils.Println("收到消息", string(msg.Payload()))
		switch string(msg.Payload()) {
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
	})
	<-closeChan
	client.Disconnect(0)
	utils.Println("done mqtt client")
}
