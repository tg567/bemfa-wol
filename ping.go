package main

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type pingResult struct {
	ip   string
	live bool
}

type pingHandler struct {
	results chan pingResult
	uid     string
	ipMap   map[string]Device
	con     net.Conn
	client  mqtt.Client
}

func newMQTTPingHandler(uid string, ipMap map[string]Device, client mqtt.Client) *pingHandler {
	return &pingHandler{
		results: make(chan pingResult, 10),
		uid:     uid,
		ipMap:   ipMap,
		client:  client,
	}
}

func newTCPPingHandler(uid string, ipMap map[string]Device, con net.Conn) *pingHandler {
	return &pingHandler{
		results: make(chan pingResult, 10),
		uid:     uid,
		ipMap:   ipMap,
		con:     con,
	}
}

func (p *pingHandler) ping(ctx context.Context, tcp bool) {
	if len(p.ipMap) == 0 {
		return
	}
	t := time.NewTicker(time.Minute * 2)
	go p.handleResult(ctx, tcp)
	for {
		select {
		case <-t.C:
			for ip, d := range p.ipMap {
				conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, strconv.Itoa(d.SSHPort)), 1*time.Second)
				if err == nil {
					p.results <- pingResult{ip, true}
					conn.Close()
				} else {
					p.results <- pingResult{ip, false}
				}
			}
		case <-ctx.Done():
			println("done ping")
			return
		}
	}
}

func (p *pingHandler) handleResult(ctx context.Context, tcp bool) {
	for {
		select {
		case result, ok := <-p.results:
			if !ok {
				return
			}
			live := `off`
			if result.live {
				live = `on`
			}
			if tcp {
				tcpMsg := fmt.Sprintf(`cmd=2&uid=%s&topic=%s/up&msg=%s`, p.uid, p.ipMap[result.ip].Topic, live)
				if _, err := p.con.Write([]byte(tcpMsg)); err != nil {
					println("更新tcp状态错误", err)
				}
			} else {
				mqttMsg := fmt.Sprintf(`%s/up`, p.ipMap[result.ip].Topic)
				if err := p.client.Publish(mqttMsg, 1, false, live).Error(); err != nil {
					println("更新mqtt状态错误", err)
				}
			}
		case <-ctx.Done():
			println("done ping result")
			return
		}
	}

}
