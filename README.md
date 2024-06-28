## 巴法云实现小爱同学语音控制电脑开关机

实现mqtt和tcp两种方式。

### 使用方式
需要golang环境打包，开机功能需要同一网段下运行

1.git clone项目。

2.```GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o wol```打包命令，根据使用环境更改GOOS和GOARCH的值；运行在路由器上等存储空间小的设备，打包建议用tcp连接方式，并且注释掉mqtt相关代码运行```go mod tidy```命令后再打包，减小执行文件大小。

3.运行命令：```./wol -uid xxx -topic xxx -mac 00:00:00:00:00:00 -broadcast 192.168.1.255 -ssh root@192.168.1.1```

参数：
```
  -broadcast string
        broadcast address
  -f string
        log file path
  -mac string
        mac
  -ssh string
        ssh user@ipaddress
  -topic string
        topic
  -type string
        wol type, tcp/mqtt (default "tcp")
  -uid string
        bemfa uid
```

### 免责

本项目为自用项目，对新手小白不友好