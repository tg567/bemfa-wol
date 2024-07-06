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
var LogFile, ConfigFile *os.File

func LoadConfig() {
	pwd, err := os.Getwd()
	if err != nil {
		Println("获取当前目录错误", err)
		return
	}

	ConfigFile, err = os.OpenFile(path.Join(pwd, configFileName), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		Println("打开配置文件错误", err)
		return
	}
	LogFile, err = os.OpenFile(path.Join(pwd, logFileName), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		Println("打开配置文件错误", err)
		return
	}
	if err := yaml.NewDecoder(ConfigFile).Decode(&WolConfig); err != nil {
		Println("解析配置文件错误", err)
		return
	}
}

func Println(a ...any) {
	log.Println(a...)
	if LogFile != nil {
		b := []any{time.Now().Format("2006-01-02 15:04:05")}
		b = append(b, a...)
		fmt.Fprintln(LogFile, b...)
	}
}
