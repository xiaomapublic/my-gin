#! /bin/bash

docker run --rm -v /home/www/gopath/src/my-gin:/go/src/my-gin -v /etc/localtime:/etc/localtime:ro -w /go/src/my-gin -e GOOS="linux" -e GOARCH="amd64" golang:latest go build -v -ldflags "-w -s" -o web
cd ../
./web

#将项目编译为二进制文件