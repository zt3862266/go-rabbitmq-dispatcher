#!/bin/bash

cd ..
echo "kill old process ..."
kill -INT `cat run/grd.pid`

echo "sleep 5 second ..."
sleep 5

echo "start new process ..."
nohup ./go-rabbitmq-dispatcher -c config/queue.yaml -l log/grd.log -pidfile run/grd.pid >/dev/null 2>&1 &

exit 0