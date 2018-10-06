# Hermogo

Message Service, default message encoding format is json

## Queue Message
队列模型旨在提供高可靠高并发的一对一消费模型。即队列中的每一条消息都只能够被某一个消费者进行消费。

### Send queue message
- If queue not exist, create new quene, send message again
- merge message and meta data for struct, then send to MNS queue
- Func SendQueueMessage parameter
    - queue: name of queue
    - v:     any thing that can be encoded to json
    - delay: seconds of delay, max 604800
- example
```go
        example/queue.go
```

### Handle queue message
- Support handle RAW/Struct format data
- Batch receive message, then batch handle data by limited pool
- Compatible handle data of older versions
- Support subscribe topic
```go
        example/server.go
```

## Topic Message
主题订阅模型旨在提供一对多的发布订阅以及消息通知功能，支持用户实现一站式多种消息通知方式：

- 推送到用户指定 HttpServer
- 推送到用户指定的 Queue（用户可以从该 Queue 拉取消息）
- 推送到邮件（组）
- 推送到短信（列表）
- WebSocket方式推送（即将支持）
- 移动推送（计划支持）

### Publish message
- If topic not exist, create new topic, send message again
- Func PublishMessage parameter
    - topic: name of topic
    - tag:   message tag, if subscriber set filter tag like this tag, subscriber cannot receive this message
    - v:     any thing that can be encoded to json
```go
        example/topic.go
```

## TODO List
- subscribe topic support email and phone endpoint

## References
    MNS - https://help.aliyun.com/document_detail/27414.html