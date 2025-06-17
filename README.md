## 巴法云实现小爱同学语音控制电脑开关机

实现mqtt和tcp两种方式。增加检测设备在线状态，更新巴法云设备开关状态。

开机原理是使用mac地址广播网络唤醒，关机原理是ssh连接到设备执行shutdown命令关机，检测设备在线状态原理是检查设备ssh端口是否打开。

### 打包

需要golang环境打包

#### 命令行打包

1.git clone项目,切换到main分支。

2.```GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o wol```打包命令，根据使用环境更改GOOS和GOARCH的值；运行在路由器上等存储空间小的设备，打包建议用tcp连接方式，并且注释掉mqtt相关代码运行```go mod tidy```命令后再打包，减小执行文件大小。

#### GUI打包

gui版本已经没更新了，等有时间再说。

1.git clone项目,切换到gui分支。

2.安装fyne
```
go install fyne.io/fyne/v2/cmd/fyne@latest
fyne install
```

3.```CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc fyne release -os windows -appID com.tg567.bemfawol -appVersion 1.0.0 -icon ./icon.png -appBuild 1 -developer tg567 -certificate 123456 -password 123456```打包命令,参数参考[fyne项目](https://github.com/fyne-io/fyne)或命令```fyne release -h```

### 使用

去掉传参运行，配置项较多，改为配置文件

#### 配置文件运行

运行命令：```./wol -c config.yaml```

yaml配置可参考config.yaml.example文件，配置项如下:
```yaml
uid: xxxxxxxxxxxxxxxx
log_file: ./wol.log
type: tcp
devices: 
  - mac: 00:00:00:00:00:00
    topic: xxxx
# broadcast不填时，取ip地址对应的广播地址，没有ip地址时必填
    broadcast: 192.168.1.255
    user: xxx
    ip: 192.168.1.1
#    ssh_port: 22
  - mac: 00:00:00:00:00:01
    topic: xxxxx
# broadcast不填时，取ip地址对应的广播地址，没有ip地址时必填
    broadcast: 192.168.1.255
    user: xxx
    ip: 192.168.1.2
#    ssh_port: 22
```