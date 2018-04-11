#!/bin/bash

echo "start new process ..."
nohup ./go-rabbitmq-dispatcher -c config/queue.yaml -l log/grd.log -pidfile run/grd.pid >/dev/null 2>&1 &

exit 0