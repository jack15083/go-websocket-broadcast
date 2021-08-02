## 简介

golang websocket 服务可通过http接口push消息到web客户端，消息发送采用golang的并发模式发送，并异步记录发送的消息日志。

## 安装

- 导入 db.sql 安装相关push日志表
- 更改 config.dev.json中的相关db配置与项目路径配置
- 执行 go build -o build/push_service
- 执行 build/push_service local


## 注意

本项目websocket使用用户token验证连接，这一块验证的逻辑需要根据自己的业务去更改或删除。

## 使用

例如使用PHP客户端push消息 

```php

<?php
/**
 * Created by PhpStorm.
 * Date: 2019/5/23
 * Time: 17:03
 */

require_once 'XHXPushApi.php';

$title = "通知标题pxy";
$content = "测试通知内容，测试通知内容，测试通知内容，测试通知内容，测试通知内容，测试通知内容。";

/**  push接口支持的参数
     * @param $senderId      //发送者id
     * @param $senderName    //发送者姓名
     * @param $msgType       //消息类型 1发送在线用户即时消息 2登录后必达消息 3业务内容更新消息
     * @param $title         //消息标题
     * @param $content       //消息内容数组或字符串，如果是数组将会被json_encode
     * @param $userIds       //用户id以,号分隔 msgType为2时userIds必传
     * @param $options       //弹窗选项目前支持 duration(毫秒), position, type参数（对应elementUi通知组件参数）
     * @return array
*/
$res = XHXPushApi::getInstance()->push(1,'god', 2, $title, $content, [44, 63]);

print_r($res);

```

client.html

```html
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <!-- import CSS -->
  <link rel="stylesheet" href="https://unpkg.com/element-ui/lib/theme-chalk/index.css">
</head>
<body>
  <div id="app">
    
  </div>
</body>
  <!-- import Vue before Element -->
  <script src="https://unpkg.com/vue/dist/vue.js"></script>
  <!-- import JavaScript -->
  <script src="https://unpkg.com/element-ui/lib/index.js"></script>
  <script>
    new Vue({
      el: '#app',
      data() {
        return { 
            visible: false,
            ws:'',
            interval:'',
            retryConnect:false,
        }
      },
      created() {
          this.init()
      },
      methods: {
        init() {
            if (!window["WebSocket"]) {
                console.log('not support websocket')
                return
            }

            var that = this;
            this.ws = new WebSocket('ws://127.0.0.1:9002/ws/' + '?token=xxxxx&uid=xxx');
            this.ws.onclose = function(e) {
                clearInterval(that.interval)
                if(!that.retryConnect) {
                    return
                }
                console.log('push connection is close, retry connect after 5 seconds')
                setTimeout(function() {
                    that.init()
                }, 5000);
            }
            this.ws.addEventListener('open', function (e) {
                
            });

            this.ws.addEventListener("message", function(e) {
                let res = JSON.parse(e.data)
                
                //token过期
                if(res.error == 100) {
                    console.log(res)
                    that.retryConnect = false
                    return
                }

                if(res.error != 0) {
                    console.log(res.msg)
                    return
                }
                
                //client注册消息
                if(res.event == 'register') {
                    console.log('ws connection register success ')
                    that.interval = setInterval(function() {
                        //保此常连接心跳
                        that.ws.send('{}')
                    }, 60000)
                    that.retryConnect = true
                    return;
                }

                if(res.event == 'message') {
                    let options = JSON.parse(res.data.options);
                    that.$notify.info({
                        title: res.data.title != '' ? res.data.title : '通知',
                        message: res.data.content,
                        duration: options.duration,
                        position: options.position
                    });
                }
            })
        }
      }
    })
  </script>
</html>
```
