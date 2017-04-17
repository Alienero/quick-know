# Quick-Know 中文文档

Quick-Know是一个高可用的推送集群服务。
可创建多个App级用户，其下可以创建多个Client终端用户。Client级用户的消息由App级用户进行管理和消息推送。

## 特性
- 易部署
- 使用Etcd做故障转移及配置文件的分享
- 多数据库支持
- 默认MongoDB提供快速的离线消息存储
- 基于Mqtt协议的推送
- 多个App用户
- 每个App下可以拥有多个Client
- 支持对App内所用用户广播消息
- 支持对App内用户私信推送
- 支持App内添加多个订阅组（类似IM聊天系统的群）
- 支持App内消息过期
- 支持Tcp推送与Websocket推送
- 支持离线消息存储
- 应用层心跳，保证用户在线可靠性
- 完善的Restful API，为用户提供全面的对App操作

## 架构
![quick-know](https://raw.githubusercontent.com/Alienero/quick-know/master/Documentation/img/qk.png "Quick-know")

## 快速开始
[中文](https://github.com/Alienero/quick-know/blob/master/Documentation/quick_start.md)