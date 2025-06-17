package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"gopkg.in/yaml.v3"
)

type Setting struct {
	UID     string   `yaml:"uid"`
	Devices []Device `yaml:"devices"`
	LogFile string   `yaml:"log_file"`
	Type    string   `yaml:"type"`
}

type Device struct {
	Topic     string `yaml:"topic"`     //"xxx"
	MAC       string `yaml:"mac"`       //"00:00:00:00:00:00"
	Broadcast string `yaml:"broadcast"` //"192.168.1.255"
	User      string `yaml:"user"`
	IP        string `yaml:"ip"`
	SSHPort   int    `yaml:"ssh_port"`
}

func (s *Setting) Validate() error {
	if s.UID == "" {
		return errors.New("UID is empty")
	}
	for i, v := range s.Devices {
		if v.Topic == "" {
			return errors.New("Topic is empty")
		}
		if v.MAC == "" {
			return errors.New("MAC is empty")
		}
		if v.SSHPort == 0 {
			s.Devices[i].SSHPort = 22
		}
		if v.Broadcast == "" {
			broadcast := generateBroadCast(v.IP)
			if broadcast == "" {
				return errors.New("Broadcast is empty")
			}
			s.Devices[i].Broadcast = broadcast
		}
	}
	return nil
}

func loadSetting(configPath string) (*Setting, error) {
	var setting Setting
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bs, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(bs, &setting); err != nil {
		return nil, err
	}
	return &setting, nil
}

func generateBroadCast(ip string) string {
	targetIP := net.ParseIP(ip)
	if targetIP == nil {
		println("无效的IP地址")
		return ""
	}

	mask := targetIP.DefaultMask()
	if mask == nil {
		fmt.Println("无法获取子网掩码")
		return ""
	}

	network := targetIP.Mask(mask)

	broadcast := make(net.IP, len(network))
	for i := range network {
		broadcast[i] = network[i] | ^mask[i]
	}
	return broadcast.String()
}
