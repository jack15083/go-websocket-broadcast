package models

import (
	"time"

	"push_service/src/config"
)

type PushMessageModel struct {
	ID         int64
	Title      string
	Content    string
	Options    string
	MsgType    int
	UserIds    string
	SenderId   int64
	SenderName string
	CreateTime string
}

func (PushMessageModel) TableName() string {
	return "xhx_push_message"
}

func (PushMessageModel) Create(m PushMessageModel) int64 {
	db := GetDB("default")
	m.CreateTime = time.Now().Format(config.TIMESTAMP_FORMAT)
	db.Create(&m)

	return m.ID
}
