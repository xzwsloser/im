# 即时通信系统(im)

> Go-Zero + Mysql + Redis + MongoDB + Kafka

这是一个微服务架构的 `im` 系统,包含用户服务,社交服务,`im` 服务以及消息队列服务,这里的核心在于 `im`服务以及消息队列服务,其他的服务只是涉及到简单的 `CRUD` 业务



目前 `im` 服务以及 消息队列服务的完成的基本功能如下:

- 构建 `websocket server` , 支持用户登录到 `im` 服务
- 用户在线检测以及`websocket`连接心跳检测
- 好友私聊(使用 `Kafka` 进行异步消息处理)
- 消息确认机制 `ACK` , 保障消息可靠性,减少消息丢失
- 离线消息拉取
- 群组聊天
- 利用 `bitmap` 实现消息已读未读功能
- 利用`Redis`存储用户状态,检测用户状态(在线,离线)



整个 `im` 系统的消息的收发模型如下:

![image-20250323180129632](/home/loser/.config/Typora/typora-user-images/image-20250323180129632.png)







