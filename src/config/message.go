package config

//客户端注册消息
type RegisterMessage struct {
	Token string
	Event string
}

//服务端返回消息
type ResMessage struct {
	Error int         `json:"error"`
	Msg   string      `json:"msg"`
	Data  MessageData `json:"data"`
	Event string      `json:"event"`
}

//推送数据结构
type PushMessage struct {
	SenderId   int64  //发送者id
	SenderName string //发送者姓名
	MsgType    int    //消息类型 1发送在线用户即时消息 2登录后必达消息 3业务内容更新消息
	Title      string //消息标题
	Content    string //消息内容
	UserIds    string //用户id以,号分隔 msgType为2时userIds必传
	Options    string //弹窗选项目前支持 duration(毫秒), position, type参数（对应elementUi通知组件参数）
	Timestamp  string //时间戳
}

//发送给客端data消息数据结构
type MessageData struct {
	SenderId   int64  `json:"senderId"`
	SenderName string `json:"senderName"`
	MsgTime    string `json:"msgTime"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Options    string `json:"options"`
	MsgId      int64  `json:"msgId,omitempty"`    //推送消息数据库记录id
	MsgType    int    `json:"msgType,omitempty"`  //消息类型
	MsgLogId   int64  `json:"msgLogId,omitempty"` //用户消息数据库记录id
}

//消息选项数据结构
type MessageOptions struct {
	Duration int    `json:"duration"`
	Position string `json:"position"`
	Type     string `json:"type"`
	PopType  string `json:"popType"`
}
