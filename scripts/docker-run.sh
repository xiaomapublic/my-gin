#!/bin/bash
docker run -it --name go-mygin -p 8181:8181 -v /home/www/gopath/src/my-gin:/go/src/my-gin -v /home/docker/gomygin:/logs -v /etc/localtime:/etc/localtime:ro -w /go/src/my-gin be63d15101cb go run /go/src/my-gin

#运行容器，运行项目时使用