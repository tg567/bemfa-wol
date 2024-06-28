## 巴法云实现小爱同学语音控制电脑开关机

实现mqtt和tcp两种方式，自用tcp方式，mqtt方式已实现，但是注释了。

### 使用方式
需要golang环境打包，开机功能需要同一网段下运行

1.git clone项目，修改main.go中uid，topic，mac，ipAddress，shutdownUserServer为自己环境的值。

2.```GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o wol```打包命令，根据使用环境更改GOOS和GOARCH的值。

3.```./wol -f ./wol.log```运行项目，-f 参数为日志文件路径，可省略仅输出到命令行。

### 免责

本项目为自用项目，对新手小白不友好