package utils

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	UID       string `yaml:"uid"`
	Topic     string `yaml:"topic"`
	Mac       string `yaml:"mac"`
	Broadcast string `yaml:"broadcast"`
	SSH       string `yaml:"ssh"`
	Type      string `yaml:"type"`
}

var WolConfig Config

var logFileName, configFileName = "wol.log", "wol.yaml"
var logFile, configFile *os.File

const (
	WOL_TYPE_TCP  = "TCP"
	WOL_TYPE_MQTT = "MQTT"
)

func LoadConfig() error {
	pwd, err := os.Getwd()
	if err != nil {
		Println("获取当前目录错误", err)
		return err
	}

	configFile, err = os.OpenFile(path.Join(pwd, configFileName), os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		Println("打开配置文件错误", err)
		return err
	}
	logFile, err = os.OpenFile(path.Join(pwd, logFileName), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		Println("打开配置文件错误", err)
		return err
	}
	if err := yaml.NewDecoder(configFile).Decode(&WolConfig); err != nil {
		Println("解析配置文件错误", err)
	}
	return nil
}

func SaveConfig() error {
	if err := configFile.Truncate(0); err != nil {
		return err
	}
	if _, err := configFile.Seek(0, 0); err != nil {
		return err
	}
	return yaml.NewEncoder(configFile).Encode(WolConfig)
}

func Println(a ...any) {
	log.Println(a...)
	if logFile != nil {
		b := []any{time.Now().Format("2006-01-02 15:04:05")}
		b = append(b, a...)
		fmt.Fprintln(logFile, b...)
	}
}
