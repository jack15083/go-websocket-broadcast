package models

import (
	"time"

	"../config"
)

type PushMessageLogModel struct {
	ID         int64
	MsgId      int64
	MsgType    int
	ClientId   string
	UserId     int64
	Status     int
	Deleted    int
	CreateTime string
	UpdateTime string
}

func (PushMessageLogModel) TableName() string {
	return "xhx_push_message_log"
}

//写入待发送消息日志
func (PushMessageLogModel) CreateWaiteMessageLogs(waitUserIds []interface{}, msgId int64, msgType int, createTime string) {
	if msgId == 0 {
		return
	}

	db, err := BaseModel.ConnectDB("default")
	if err != nil {
		return
	}

	for _, waitUserId := range waitUserIds {
		userId := waitUserId.(int64)
		if userId <= 0 {
			continue
		}
		pml := PushMessageLogModel{MsgId: msgId, MsgType: msgType, CreateTime: createTime, UserId: userId, UpdateTime: createTime}
		db.Create(&pml)
	}

	defer db.Close()
}

//新增发送日志 status 1 发送成功 2发送失败 delete是否册除
func (PushMessageLogModel) Create(msgId int64, msgType int, userId int64, clientId string, status int) int64 {
	db, err := BaseModel.ConnectDB("default")
	if err != nil {
		return 0
	}
	defer db.Close()

	if msgId == 0 {
		return 0
	}

	nowTime := time.Now().Format(config.TIMESTAMP_FORMAT)
	m := &PushMessageLogModel{MsgId: msgId, MsgType: msgType, ClientId: clientId, Status: status, UserId: userId,
		CreateTime: nowTime, UpdateTime: nowTime}

	db.Create(m)

	return m.ID
}

//更新消息
func (pml PushMessageLogModel) Save(id int64, clientId string, status int) {
	db, err := BaseModel.ConnectDB("default")
	if err != nil {
		return
	}
	defer db.Close()
	nowTime := time.Now().Format(config.TIMESTAMP_FORMAT)
	db.Model(&pml).Where("id = ?", id).Updates(PushMessageLogModel{ClientId: clientId, Status: status, UpdateTime: nowTime})
}

func (pml PushMessageLogModel) SetDeleted(id int64) {
	db, err := BaseModel.ConnectDB("default")
	if err != nil {
		return
	}
	defer db.Close()
	nowTime := time.Now().Format(config.TIMESTAMP_FORMAT)
	db.Model(&pml).Where("id = ?", id).Updates(PushMessageLogModel{Deleted: 1, UpdateTime: nowTime})
}

//获取用户最近的必读消息
func (pml PushMessageLogModel) GetMustReadMsgByUserId(userId int64, unixtime int64) []config.MessageData {
	db, err := BaseModel.ConnectDB("default")
	var mustReadData []config.MessageData
	if err != nil {
		return mustReadData
	}
	defer db.Close()

	var data []PushMessageLogModel
	db.Where("user_id = ? and msg_type = 2 and status in(0, 2) and deleted = 0 and UNIX_TIMESTAMP(create_time) > ? ",
		userId, unixtime).Limit(config.LAST_MSG_NUM_LIMIT).Find(&data)

	for _, row := range data {
		//对发送失败的消息做重复发送检查
		if row.Status == 2 {
			count := 0
			db.Model(&pml).Where("user_id = ? and msg_id = ? and status = 1 and deleted = 0", userId, row.MsgId).Count(&count)
			if count > 0 {
				nowTime := time.Now().Format(config.TIMESTAMP_FORMAT)
				//删除重复的消息
				db.Model(&pml).Where("id = ?", row.ID).Updates(PushMessageLogModel{Deleted: 1, UpdateTime: nowTime})
				continue
			}
		}

		var pm PushMessageModel
		db.First(&pm, row.MsgId)

		mustReadData = append(mustReadData, config.MessageData{SenderId: pm.SenderId, SenderName: pm.SenderName, MsgTime: pm.CreateTime, Title: pm.Title,
			Content: pm.Content, Options: pm.Options, MsgId: pm.ID, MsgType: pm.MsgType, MsgLogId: row.ID})
	}

	return mustReadData
}
