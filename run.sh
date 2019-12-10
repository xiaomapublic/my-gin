#!/usr/bin/env bash
#git pull

go build -o my-gin-linux-amd64 -ldflags="-s -w" -tags=jsoniter -v ./
echo "my-gin-linux-amd64 build success"

pid=`ps -ef|grep -v grep|grep ./my-gin-linux-amd64 |awk '{print $2}'`
if [ -n "$pid" ]
then
    echo "kill pid:" $pid
    kill $pid
fi

sleep 5s

nohup ./my-gin-linux-amd64 &

#pid=`ps -ef|grep -v grep|grep ./my-gin-linux-amd64 |awk '{print $2}'`

#echo "run pid:" $pid

#sleep 2s
#tail -f nohup.out
