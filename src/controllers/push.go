package controllers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"push_service/src/config"
	"push_service/src/models"
	"push_service/src/server"
	"push_service/src/utils"
)

type PushController struct {
	BaseController
}

func (c *PushController) Hello(w http.ResponseWriter, r *http.Request) {
	count := 0
	server.Manager.Clients.Range(func(k, v interface{}) bool {
		count++
		return true
	})

	w.Write([]byte(fmt.Sprintf("当前client连接总数：%d\n\n", count)))
	w.Write([]byte("client连接详情：\n"))

	server.Manager.Clients.Range(func(k, v interface{}) bool {
		conn := k.(*server.Client)
		w.Write([]byte(fmt.Sprintf("Client ID:%s, Admin ID:%d\n", conn.ID, conn.UserId)))
		return true
	})
}

func (c *PushController) Push(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method not allowed!")
		return
	}

	// read request
	var pm config.PushMessage
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&pm); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error("Push message error: " + err.Error())
		return
	}

	if len(r.Header["Access-Token"]) == 0 {
		c.sendError(w, 400, "AccessToken校验失败")
	}

	if !c.checkApiKey(r.Header["Access-Token"][0], pm) {
		c.sendError(w, 401, "AccessToken验证失败")
		return
	}

	if !c.checkPushParams(w, pm) {
		return
	}

	//数据写入到数据库
	var pushMsgModel models.PushMessageModel
	msgId := pushMsgModel.Create(models.PushMessageModel{SenderId: pm.SenderId, SenderName: pm.SenderName, Title: pm.Title, Content: pm.Content,
		Options: pm.Options, MsgType: pm.MsgType, UserIds: pm.UserIds})

	data := config.MessageData{SenderId: pm.SenderId, MsgTime: time.Now().Format(config.TIMESTAMP_FORMAT), SenderName: pm.SenderName, Title: pm.Title, Content: pm.Content,
		Options: pm.Options, MsgId: msgId, MsgType: pm.MsgType}
	message, _ := json.Marshal(&config.ResMessage{Error: 0, Msg: "ok", Event: "message", Data: data})

	if pm.UserIds == "0" { //发全部
		server.Manager.Broadcast <- message
	} else {
		userIdsArr := strings.Split(pm.UserIds, ",")
		var userIds = make([]interface{}, 0)
		for _, userId := range userIdsArr {
			userId = strings.Trim(userId, " ")
			userId, _ := strconv.ParseInt(userId, 10, 64)
			userIds = append(userIds, userId)
		}

		if len(userIds) > config.MAX_SEND_USER_NUM && pm.MsgType == 2 {
			c.sendError(w, 200, fmt.Sprintf("必读消指定用户时最多发送用户量不可超过%d", config.MAX_SEND_USER_NUM))
			return
		}
		//发送消息到指定用户
		sendUserIds := server.Manager.SendMsgToUsers(message, userIds)
		if pm.MsgType == 2 { //如果是必达消息并发模式写入数据库状态为待发送
			go func() {
				waitSendUserIds := utils.SliceDiff(userIds, sendUserIds)
				var pushMsgLogModel models.PushMessageLogModel
				pushMsgLogModel.CreateWaiteMessageLogs(waitSendUserIds, msgId, pm.MsgType, time.Now().Format(config.TIMESTAMP_FORMAT))
			}()
		}
	}

	c.sendOk(w, "ok")
}

//检查push接口参数
func (c *PushController) checkPushParams(w http.ResponseWriter, pm config.PushMessage) bool {
	msgTime, err := strconv.ParseFloat(pm.Timestamp, 64)
	if err != nil {
		fmt.Println(err)
	}

	if math.Abs(float64(time.Now().Unix())-msgTime) > 120 {
		c.sendError(w, 109, "Token expired")
		return false
	}

	if pm.MsgType != 1 && pm.MsgType != 2 && pm.MsgType != 3 {
		c.sendError(w, 102, "MsgType参数错误")
		return false
	}

	if pm.UserIds == "" {
		c.sendError(w, 103, "UserIds必传")
		return false
	}

	if pm.SenderId == 0 {
		c.sendError(w, 104, "SenderId必传")
		return false
	}

	if pm.SenderName == "" {
		c.sendError(w, 105, "SenderName必传")
		return false
	}

	if pm.Content == "" {
		c.sendError(w, 106, "Content必传")
		return false
	}

	if pm.MsgType == 2 && pm.UserIds == "0" {
		c.sendError(w, 108, "msgType为2时必须传要发送的用户id")
		return false
	}

	return true
}
