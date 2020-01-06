#!/bin/bash

APP_ENV=$1
APP_ROOT_PATH="/data/app/xhx_pushservice/src"
APP_DIST_PATH="/data/app/xhx_pushservice/dist"
APP_NAME="xhx_pushservice"

branch=""

if [[ "$APP_ENV" == '' ]];then
    echo 'Please add env param!'
    exit
fi

go build -o build/$APP_NAME

##重启服务##
PID=`ps aux | grep $APP_NAME.*.$APP_ENV | grep -v grep  | awk '{ print $2}'`
if [[ "" !=  "$PID" ]];then
    echo "killing old service[$PID]"
    kill -9 $PID
fi

nohup build/$APP_NAME -c config.$APP_ENV.json >/dev/null 2>&1 &

echo -e "$APP_NAME start \033[32;40m [SUCCESS] \033[0m"