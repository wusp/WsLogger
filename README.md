# WsLogger
A command line tool used to record WebSocket stream into log files.

## Download Binary

| OS | Link |
| --- | --- |
| Android | https://wusp-file.oss-cn-shenzhen.aliyuncs.com/android-WsLogger
| Linux | https://wusp-file.oss-cn-shenzhen.aliyuncs.com/linux-WsLogger |
| Windows | https://wusp-file.oss-cn-shenzhen.aliyuncs.com/win-WsLogger.exe |
| Mac | https://wusp-file.oss-cn-shenzhen.aliyuncs.com/mac-WsLogger |


## Run

    ./WsLogger -a localhost:8080

Record websocket stream which is running on localhost:8080  into ${PWD}/log dir.

##### Options
| Args | Description | Default Value |
| --- | --- | --- |
| -h | Tool Help | - |
| -a | Connect Url | localhost:8081 |
| -c | Auto Reconnect, false if -c is absent. | false |
| -d | Log Store Dir  | $(PWD)/log |
| -p | Url Path | /rawcomm |
| -t | Read Timeout in second| 60 |


## Build


    make dist
    
4 different OS productions would be appeared on ${rootProject}/dist/.


