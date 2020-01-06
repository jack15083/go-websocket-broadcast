/*
SQLyog Professional v12.08 (64 bit)
MySQL - 5.7.25-log : Database - push_service
*********************************************************************
*/


/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`push_service` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

USE `push_service`;

/*Table structure for table `xhx_push_message` */

DROP TABLE IF EXISTS `xhx_push_message`;

CREATE TABLE `xhx_push_message` (
  `id` bigint(15) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) DEFAULT '' COMMENT '消息业务标题',
  `content` text NOT NULL COMMENT '消息业务内容',
  `options` text COMMENT '消息弹窗配置',
  `msg_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '消息类型 1发送在线用户即时消息 2登录后必达消息 3 业务内容更新消息',
  `user_ids` text NOT NULL COMMENT '要发送的用户id 0表示发全部',
  `sender_id` bigint(15) NOT NULL COMMENT '发送者id',
  `sender_name` varchar(100) NOT NULL COMMENT '发送者姓名',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_sid` (`sender_id`)
) ENGINE=InnoDB AUTO_INCREMENT=440 DEFAULT CHARSET=utf8mb4;

/*Table structure for table `xhx_push_message_log` */

DROP TABLE IF EXISTS `xhx_push_message_log`;

CREATE TABLE `xhx_push_message_log` (
  `id` bigint(15) NOT NULL AUTO_INCREMENT,
  `msg_id` bigint(15) NOT NULL COMMENT 'push message表消息id',
  `msg_type` tinyint(1) NOT NULL COMMENT '消息类型1即时消息 2必达',
  `client_id` varchar(60) DEFAULT '' COMMENT '客户端id',
  `user_id` bigint(15) NOT NULL COMMENT '接收者用户id',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0待发送 1发送成功 2发送失败',
  `deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除0未删1已删除',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_uid` (`user_id`),
  KEY `idx_mid` (`msg_id`)
) ENGINE=InnoDB AUTO_INCREMENT=720 DEFAULT CHARSET=utf8mb4;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
