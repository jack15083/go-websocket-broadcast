<?php
/**
 * Created by PhpStorm.
 * User: zengfanwei
 * Date: 2019/5/23
 * Time: 16:26
 */

class XHXPushApi
{
    private static $_api_secret = 'dev'; //根据环境配置
    private static $_api_url	= 'http://127.0.0.1:9003/push'; //根据环境配置

     /**
     * singleton instance
     *
     */
    protected static $_instance = null;

    public function __construct()
    {
        $this->apiSecret = Yaf_Registry::get('config')->pushserver->secret;
        $this->apiUrl = Yaf_Registry::get('config')->pushserver->apiurl;
    }

    /**
     * Returns singleton instance of XHXPushApi
     *
     * @return XHXPushApi
     */
    public static function getInstance()
    {
        if (!isset( self::$_instance ))
        {
            self::$_instance = new self();
        }

        return self::$_instance;
    }


    /**
     * @param $senderId      //发送者id
     * @param $senderName    //发送者姓名
     * @param $msgType       //消息类型 1发送在线用户即时消息 2登录后必达消息 3业务内容更新消息
     * @param $title         //消息标题
     * @param $content       //消息内容数组或字符串，如果是数组将会被json_encode
     * @param $userIds       //用户id以,号分隔 msgType为2时userIds必传
     * @param $options       //弹窗选项目前支持 duration(毫秒), position, type参数（对应elementUi通知组件参数）
     *                       //popType 弹出类型 ele为elementUi通知，browser为浏览器通知，all为elementUi通知 + 浏览器通知
     * @return array
     */
    public function push($senderId, $senderName, $msgType, $title, $content, $userIds = [0], $options = [])
    {
        if(empty($options) && $msgType != 3) {
            $options = [
                'duration' => 0,
                'position' => 'top-right',
                'type' => 'info',
                'popType' => 'all'
            ];
        }

        if($msgType != 3) $options['duration'] = (int) $options['duration'];

        $params = [
            'senderId'   => (int) $senderId,
            'senderName' => $senderName,
            'msgType'    => (int) $msgType,
            'title'      => $title,
            'content'    => $content,
            'userIds'    => implode(',', $userIds),
            'options'    => json_encode($options),
            'timestamp'  => (string) microtime(true)
        ];

        $dataJson = json_encode($params);

        $token = md5(implode($this->apiSecret, $params));
        $ch = curl_init();

        curl_setopt($ch, CURLOPT_URL, $this->apiUrl);
        curl_setopt($ch, CURLOPT_HTTPHEADER, [
            'Content-Type: application/json','Content-Length: ' . strlen($dataJson),
            'ACCESS-TOKEN:' . $token
        ]);
        curl_setopt($ch, CURLOPT_CUSTOMREQUEST, 'PUT');
        curl_setopt($ch, CURLOPT_POSTFIELDS,$dataJson);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_TIMEOUT, 5);

        $data = curl_exec($ch);

        $res[0] = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        $res[1] = $data;
        if($res[0] != 200){
            $res[2] = curl_error($ch);
        }

        return $res;
    }
}