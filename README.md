## 巴法云实现小爱同学语音控制电脑开关机

实现mqtt和tcp两种方式。

### 打包

需要golang环境打包

#### 命令行打包

1.git clone项目,切换到main分支。

2.```GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o wol```打包命令，根据使用环境更改GOOS和GOARCH的值；运行在路由器上等存储空间小的设备，打包建议用tcp连接方式，并且注释掉mqtt相关代码运行```go mod tidy```命令后再打包，减小执行文件大小。

#### GUI打包

1.git clone项目,切换到gui分支。

2.安装fyne
```
go install fyne.io/fyne/v2/cmd/fyne@latest
fyne install
```

3.```CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc fyne release -os windows -appID com.tg567.bemfawol -appVersion 1.0.0 -icon ./icon.png -appBuild 1 -developer tg567 -certificate 123456 -password 123456```打包命令,参数参考[fyne项目](https://github.com/fyne-io/fyne)或命令```fyne release -h```

### 使用

开机功能需要同一网段下运行，命令行传参只支持单个设备唤醒开机关机，配置文件支持多个设备唤醒开机关机

#### 命令行传参

运行命令：```./wol -uid xxx -topic xxx -mac 00:00:00:00:00:00 -broadcast 192.168.1.255 -ssh root@192.168.1.1```

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

#### 配置文件运行

运行命令：```./wol -c config.yaml```

yaml配置可参考config.yaml.example文件，配置项如下:
```yaml
uid: xxxxxxxxxxxxxxxxxxxxx
log_file: ./wol.log
type: tcp
devices: 
  - mac: 00:00:00:00:00:00
    topic: xxx1
    broadcast: 192.168.1.255
    ssh: root@192.168.1.1
  - mac: 00:00:00:00:00:01
    topic: xxx2
    broadcast: 192.168.1.255
    ssh: root@192.168.1.2
```