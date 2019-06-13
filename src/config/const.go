package config

const XHX_SESSION_NAME = "php_xhx_token:"      //token session name
const CLIENT_TIMEOUT = 43200                   //客户端连接最大时间
const MAX_SEND_USER_NUM = 1000                 //必达消息最大发送用户量
const LAST_MSG_TIME_LIMIT = 3 * 24 * 3600      //读取最近必达消息时间限制
const LAST_MSG_NUM_LIMIT = 20                  //读取最近必达消息数量限制
const TIMESTAMP_FORMAT = "2006-01-02 15:04:05" //日期格式
const CLIENT_REGISTER_TIMEOUT = 5              //客户端注册登录超时时间 5秒
