#!/bin/bash

echo "start new process ..."
ps=`ps -ef | grep go-rabbitmq-dispatcher | grep -v grep -c`
if [ $ps -gt 0 ];then
    echo "already have $ps process,give up"!
    exit 0;
fi
nohup ./go-rabbitmq-dispatcher -c config/queue.yaml -l log/grd.log -pidfile run/grd.pid >/dev/null 2>&1 &

exit 0