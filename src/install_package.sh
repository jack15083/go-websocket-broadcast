#!/bin/bash
go get -u github.com/go-redis/redis
go get -u github.com/sirupsen/logrus
go get -u github.com/jinzhu/gorm
go get -u github.com/jinzhu/gorm/dialects/mysql
go get -u github.com/gorilla/websocket
go get -u github.com/satori/go.uuid
go get -u github.com/lestrrat-go/file-rotatelogs
go get -u github.com/pkg/errors
go get -u github.com/rifflock/lfshook

echo "执行完毕"
