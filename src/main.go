package main

import (
	"io"
	"net/http"
	"time"

	"./controllers"
	"./core"
	"./server"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	//配置初始化
	core.Config.Init()
	//日志初始化
	core.Config.Logger.Init()
	//客户端管理服务启动
	go server.Manager.Start()
	//关闭超时连接客户端
	go server.Manager.CloseTask()

	log.Infof("Server started %s ...", core.Config.Listen)
	http.HandleFunc("/ws", wsPage)
	var pc controllers.PushController
	http.HandleFunc("/push", pc.Push)
	http.HandleFunc("/hello", pc.Hello)

	log.Fatal(http.ListenAndServe(core.Config.Listen, nil))
}

//广播客户端连接handle
func wsPage(res http.ResponseWriter, req *http.Request) {
	//解析一个连接
	conn, error := upgrader.Upgrade(res, req, nil)
	if error != nil {
		io.WriteString(res, "这是一个websocket.")
		return
	}

	uid, _ := uuid.NewV4()
	sha1 := uid.String()

	//初始化一个客户端对象
	client := &server.Client{ID: sha1, Socket: conn, Send: make(chan []byte), RegisterTime: time.Now().Unix()}
	//注册一个对象到channel
	server.Manager.Register <- client

	go client.Read()
	go client.Write()
}
