package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register: //新客户端加入
			manager.clients[conn] = true
			//jsonMessage, _ := json.Marshal(&Message{Content: " a new socket has connected."})
			//manager.send(jsonMessage, conn) //调用发送
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				//jsonMessage, _ := json.Marshal(&Message{Content: "a socket has disconnected."})
				//manager.send(jsonMessage, conn)
			}
		case message := <-manager.broadcast: //读到广播管道数据后的处理
			fmt.Println("读到广播管道数据: " + string(message))
			for conn := range manager.clients {
				fmt.Println("客户端id:", conn.id)

				select {
				case conn.send <- message: //调用发送给全体客户端
				default:
					fmt.Println("执行关闭连接")
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.send <- message //发送的数据写入所有的 websocket 连接 管道
		}
	}
}

//客户端写入后 激活这里读取, 客户端广播数据
func (c *Client) read() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
		fmt.Println("读关闭")
	}()

	for {

		_, message, err := c.socket.ReadMessage()
		fmt.Println("读取客户端发送的数据")
		if err != nil {
			manager.unregister <- c
			c.socket.Close()
			fmt.Println("读不到数据就关闭？")
			break
		}
		jsonMessage, _ := json.Marshal(&Message{Sender: c.id, Content: string(message)})
		manager.broadcast <- jsonMessage //激活start 程序 入广播管道
		fmt.Println("发送数据到广播")
	}
}

//写入管道后激活这个进程
func (c *Client) write() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
		fmt.Println("写关闭了")
	}()

	for {
		select {
		case message, ok := <-c.send: //这个管道有了数据 写这个消息出去
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				fmt.Println("发送关闭提示")
				return
			}

			err := c.socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				manager.unregister <- c
				c.socket.Close()
				fmt.Println("写不成功数据就关闭了")
				break
			}
			fmt.Println("写数据")
		}
	}
}

func main() {
	fmt.Println("Starting application...")
	go manager.start()
	go manager.getRedisData()
	http.HandleFunc("/ws", wsPage)
	http.ListenAndServe(":9002", nil)
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	//解析一个连接
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		//http.NotFound(res, req)
		//http 请求一个输出
		io.WriteString(res, "这是一个websocket.")
		return
	}

	uid, _ := uuid.NewV4()
	sha1 := uid.String()

	//初始化一个客户端对象
	client := &Client{id: sha1, socket: conn, send: make(chan []byte)}
	//把这个对象发送给 管道
	manager.register <- client

	//go client.read()
	go client.write()
}

//redis 消息订阅
func (manager *ClientManager) getRedisData() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "192.168.126.128:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	redisSubscript := redisClient.Subscribe("broadcast_channel")
	for {
		msg, err := redisSubscript.ReceiveMessage()
		if err != nil {
			fmt.Println("没有读取到订阅数据")
			redisClient.Close()
		}

		fmt.Println("重新读数据")
		jsonMessage, _ := json.Marshal(&Message{Sender: "admin", Content: msg.Payload})
		fmt.Println(msg.String())
		manager.broadcast <- jsonMessage //广播数据
	}
}
