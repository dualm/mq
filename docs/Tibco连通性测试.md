# Tibco连通性测试

1. 开始-运行-cmd

2. 命令行中输入` tibrvlisten -service 14450 -network ";225.14.14.5" -daemon "tcp:10.11.111.30:7500" subject test`

3. 新建命令行窗口，输入` tibrvsend -service 14450 -network ";225.14.14.5" -daemon "tcp:10.11.111.30:7500" test test`

4. 查看tibrvlisten的输出，是否有内容为`subject=test, message={DATA="test"}`的消息日志输出 


