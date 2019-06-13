#!/bin/bash

APP_ENV=$1
APP_ROOT_PATH=""
APP_NAME="push_service_$APP_ENV"

if [[ "$APP_ENV" == '' ]];then
    echo 'Please add env param!'
    exit
fi

if [[ "$APP_ENV" == 'local' ]];then
    APP_ROOT_PATH="/mnt/hgfs/git/xhx_pushservice/src"
fi

if [[ "$APP_ENV" == 'dev' ]];then
    APP_ROOT_PATH="/home/testphp/.jenkins/workspace/xhx_pushservice_dev/src"
fi

if [[ "$APP_ENV" == 'sit' ]];then
    APP_ROOT_PATH="/home/testphp/.jenkins/workspace/xhx_pushservice_sit/src"
fi

if [[ "$APP_ROOT_PATH" == '' ]];then
    echo 'enviroment param error'
    exit
fi

cd $APP_ROOT_PATH
go build -o build/$APP_NAME

##重启服务##
PID=`ps aux | grep $APP_NAME | grep -v grep  | awk '{ print $2}'`
if [[ "" !=  "$PID" ]];then
    echo "killing old service[$PID]"
    kill -9 $PID
fi

nohup build/$APP_NAME $APP_ENV >/dev/null 2>&1 &

echo -e "$APP_NAME start \033[32;40m [SUCCESS] \033[0m"
