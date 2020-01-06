package server

import (
	"encoding/json"
	"sync"
	"time"

	"../config"
	"../utils"
	log "github.com/sirupsen/logrus"
)

type ClientManager struct {
	Clients    sync.Map
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

var Manager = ClientManager{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

//服务端启动
func (manager *ClientManager) Start() {
	for {
		select {
		case conn := <-manager.Register: //新客户端加入
			manager.Clients.Store(conn, true)
			log.Debug("New client register")
		case conn := <-manager.Unregister:
			if _, ok := manager.Clients.Load(conn); ok {
				close(conn.Send)
				manager.Clients.Delete(conn)
				log.Debug("Client unregister")
			}
		case message := <-manager.Broadcast: //读到广播管道数据后的处理
			//send all message
			manager.SendAll(message, nil)
		}
	}
}

//对ignore外所有客户端发送消息
func (manager *ClientManager) SendAll(message []byte, ignore *Client) {
	popType, options, msg := manager.GetMsgPopType(message)
	var sendUserIds = make([]interface{}, 0)
	//send all message
	manager.Clients.Range(func(k, v interface{}) bool {
		conn := k.(*Client)
		if conn != ignore {
			//发送业务更新消息
			if msg.Data.MsgType == 3 {
				conn.Send <- message
				return true
			}

			//发送element通知
			if popType == "ele" || popType == "all" {
				conn.Send <- manager.UpdateMsgPopType("ele", options, msg)
			}

			//发送浏览器通知
			if popType == "browser" || popType == "all" {
				if !utils.InSliceIface(conn.UserId, sendUserIds) {
					conn.Send <- manager.UpdateMsgPopType("browser", options, msg)
					sendUserIds = append(sendUserIds, conn.UserId)
				}
			}
		}
		return true
	})
}

//发送给指定用户消息，返回已执行发送的用户id
func (manager *ClientManager) SendMsgToUsers(message []byte, userIds []interface{}) []interface{} {
	popType, options, msg := manager.GetMsgPopType(message)
	var sendUserIds = make([]interface{}, 0)
	manager.Clients.Range(func(k, v interface{}) bool {
		conn := k.(*Client)
		for _, userId := range userIds {
			if conn.UserId == userId {
				//发送业务更新消息
				if msg.Data.MsgType == 3 {
					conn.Send <- message
					continue
				}

				//发送element通知
				if popType == "ele" || popType == "all" {
					conn.Send <- manager.UpdateMsgPopType("ele", options, msg)
				}

				//发送浏览器通知
				if popType == "browser" || popType == "all" {
					if !utils.InSliceIface(conn.UserId, sendUserIds) {
						conn.Send <- manager.UpdateMsgPopType("browser", options, msg)
					}
				}
				sendUserIds = append(sendUserIds, userId)
			}
		}
		return true
	})

	sendUserIds = utils.SliceUnique(sendUserIds)

	return sendUserIds
}

//更新消息弹窗类型
func (manager *ClientManager) UpdateMsgPopType(popType string, options config.MessageOptions, msg config.ResMessage) []byte {
	options.PopType = popType
	b, _ := json.Marshal(options)
	msg.Data.Options = string(b)

	newMsg, _ := json.Marshal(msg)
	return newMsg
}

//获取消息弹窗类型
func (manager *ClientManager) GetMsgPopType(message []byte) (string, config.MessageOptions, config.ResMessage) {
	var msg config.ResMessage
	json.Unmarshal(message, &msg)
	//获取消息弹窗类型默认all
	popType := "all"
	var options config.MessageOptions
	if msg.Data.Options != "" {
		err := json.Unmarshal([]byte(msg.Data.Options), &options)
		if err == nil && options.PopType != "" {
			popType = options.PopType
		}
	}

	return popType, options, msg
}

//定时关闭超时连接任务
func (manager *ClientManager) CloseTask() {
	ticker := time.NewTicker(time.Second * 60)
	go func() {
		for range ticker.C {
			manager.Clients.Range(func(k, v interface{}) bool {
				conn := k.(*Client)
				//超过时间关闭连接
				if time.Now().Unix()-conn.RegisterTime > config.CLIENT_TIMEOUT {
					conn.Close()
				}

				//检查token是否过期
				if conn.UserId != 0 && !conn.CheckToken() {
					conn.Close()
				}

				return true
			})
		}
	}()

}
