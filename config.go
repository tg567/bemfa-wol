package main

import (
	"errors"
	"io"
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
}

type paramDevice struct {
	UID     string
	LogFile string
	Type    string
	Device
}

func (s *Setting) Validate() error {
	if s.UID == "" {
		return errors.New("UID is empty")
	}
	for _, v := range s.Devices {
		if v.Topic == "" {
			return errors.New("Topic is empty")
		}
		if v.MAC == "" {
			return errors.New("MAC is empty")
		}
		if v.Broadcast == "" {
			return errors.New("Broadcast is empty")
		}
	}
	return nil
}

func (s *Setting) initParam(p *paramDevice) {
	s.UID = p.UID
	s.Type = p.Type
	s.LogFile = p.LogFile
	s.Devices = append(s.Devices, p.Device)
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
