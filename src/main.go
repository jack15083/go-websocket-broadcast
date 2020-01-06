package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"./controllers"
	"./core"
	"./server"
	log "github.com/sirupsen/logrus"
)

func main() {
	//配置初始化
	core.Config.Init()
	//日志初始化
	core.Config.Logger.Init()

	chExit := make(chan os.Signal)
	signal.Notify(chExit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)

	//客户端管理服务启动
	go server.Manager.Start()
	//关闭超时连接客户端
	go server.Manager.CloseTask()

	//监听http服务
	go func() {
		mux := http.NewServeMux()
		var pc controllers.PushController
		mux.HandleFunc("/push", pc.Push)
		log.Infof("Http Server started %s ...", core.Config.HttpListen)
		log.Fatal(http.ListenAndServe(core.Config.HttpListen, mux))
	}()

	//监听websocket服务
	go func() {
		mux := http.NewServeMux()
		var pc controllers.PushController
		mux.HandleFunc("/ws", server.ClientRegister)
		mux.HandleFunc("/hello", pc.Hello)
		log.Infof("Websocket Server started %s ...", core.Config.Listen)
		log.Fatal(http.ListenAndServe(core.Config.Listen, mux))
	}()

	//主进程阻塞直到有退出信号
	s := <-chExit
	log.Info("Get signal:", s)
}
