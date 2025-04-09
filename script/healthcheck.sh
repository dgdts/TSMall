#!/bin/bash
#!/bin/sh

# 检查进程是否存活
#ps -ef | grep "ts_mall_service" | grep -v "grep" | wc -l | awk '{if ($1>0) exit 0; else exit -1;}'

# 检查端口是否监听
port_cnt=`netstat -tulpn | grep 'LISTEN' | grep ':8000' | wc -l`
if [[ $port_cnt -lt 1 ]]; then
    echo "health check fail, no port listen"
    exit 1
fi

# 检查进程是否存活
proc_cnt=`ps aux | grep "ts_mall_service" | grep -v "grep" | wc -l`
if [[ $proc_cnt -lt 1 ]]; then
    echo "health check fail, no process"
    exit 2
fi

exit 0