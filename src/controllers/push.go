package controllers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"../config"
	"../core"
	"../models"
	"../server"
	"../utils"
	log "github.com/sirupsen/logrus"
)

type PushController struct {
	BaseController
}

func (c *PushController) Hello(w http.ResponseWriter, r *http.Request) {
	c.sendOk(w, "Hello World")
}

func (c *PushController) Push(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method not allowed!")
		return
	}

	// read request
	var pm config.PushMessage
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&pm); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "bad request!")
		log.Error("push message error: " + err.Error())
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
			c.sendError(w, 200, fmt.Sprintf("msgType为2时最多发送用户量不可超过%d", config.MAX_SEND_USER_NUM))
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

//push api接口token校验
func (c *PushController) checkApiToken(pm config.PushMessage) bool {
	var str string

	str += strconv.FormatInt(pm.SenderId, 10) + core.Config.APISecret
	str += pm.SenderName + core.Config.APISecret
	str += strconv.Itoa(pm.MsgType) + core.Config.APISecret
	str += pm.Title + core.Config.APISecret
	str += pm.Content + core.Config.APISecret
	str += pm.UserIds + core.Config.APISecret
	str += pm.Options + core.Config.APISecret
	str += pm.Timestamp

	//方法一
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制

	if md5str != pm.Token {
		return false
	}

	return true
}

//检查push接口参数
func (c *PushController) checkPushParams(w http.ResponseWriter, pm config.PushMessage) bool {
	if pm.Token == "" {
		c.sendError(w, 101, "Token必传")
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

	if !c.checkApiToken(pm) {
		c.sendError(w, 107, "Token 验证失败")
		return false
	}

	return true
}
