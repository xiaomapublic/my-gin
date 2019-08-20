#! /bin/bash

cd ${GOPATH}/src/my-gin
rm -rf my-gin
go build -v -o my-gin-binary-file
./my-gin-binary-file

