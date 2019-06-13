package server

import (
	"time"

	"../config"
	log "github.com/sirupsen/logrus"
)

type ClientManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

var Manager = ClientManager{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[*Client]bool),
}

func (manager *ClientManager) Start() {
	for {
		select {
		case conn := <-manager.Register: //新客户端加入
			manager.Clients[conn] = true
			log.Debug("new client register")
		case conn := <-manager.Unregister:
			if _, ok := manager.Clients[conn]; ok {
				close(conn.Send)
				delete(manager.Clients, conn)
				log.Debug("client unregister")
			}
		case message := <-manager.Broadcast: //读到广播管道数据后的处理
			for conn := range manager.Clients {
				select {
				case conn.Send <- message: //send all
				default:
					log.Debug("Execute close:" + conn.ID)
					close(conn.Send)
					delete(manager.Clients, conn)
				}
			}
		}
	}
}

func (manager *ClientManager) Send(message []byte, ignore *Client) {
	for conn := range manager.Clients {
		if conn != ignore {
			conn.Send <- message //send message to ignore
		}
	}
}

//发送消息，返回已执行发送的用户id
func (manager *ClientManager) SendMsgToUsers(message []byte, userIds []interface{}) []interface{} {
	var sendUserIds = make([]interface{}, 0)
	for conn := range manager.Clients {
		for _, userId := range userIds {
			if conn.UserId == userId {
				conn.Send <- message
				sendUserIds = append(sendUserIds, userId)
			}
		}
	}
	return sendUserIds
}

//定时关闭超时连接任务
func (manager *ClientManager) CloseTask() {
	ticker := time.NewTicker(time.Second * 60)
	go func() {
		for range ticker.C {
			for conn := range manager.Clients {
				//超过时间关闭连接
				if time.Now().Unix()-conn.RegisterTime > config.CLIENT_TIMEOUT {
					conn.Close()
				}

				//超时未登录关闭连接
				if time.Now().Unix()-conn.RegisterTime > config.CLIENT_REGISTER_TIMEOUT {
					if conn.UserId == 0 {
						conn.Close()
					}
				}
			}
		}
	}()

}
