package proto

import "encoding/xml"

// CreateQueueRequest create queue request
// see https://help.aliyun.com/document_detail/35129.html
type CreateQueueRequest struct {
	XMLName                xml.Name `xml:"Queue"`
	DelaySeconds           int      `xml:"DelaySeconds,omitempty"`           // 发送到该 Queue 的所有消息默认将以DelaySeconds参数指定的秒数延后可被消费，单位为秒。-- 0-604800秒（7天）范围内某个整数值，默认值为0
	MaximumMessageSize     int      `xml:"MaximumMessageSize,omitempty"`     // 发送到该Queue的消息体的最大长度，单位为byte。	1024(1KB)-65536（64KB）范围内的某个整数值，默认值为65536（64KB）。
	MessageRetentionPeriod int      `xml:"MessageRetentionPeriod,omitempty"` // 消息在该 Queue 中最长的存活时间，从发送到该队列开始经过此参数指定的时间后，不论消息是否被取出过都将被删除，单位为秒。	60 (1分钟)-1296000 (15 天)范围内某个整数值，默认值345600 (4 天)
	VisibilityTimeout      int      `xml:"VisibilityTimeout,omitempty"`      // 消息从该 Queue 中取出后从Active状态变成Inactive状态后的持续时间，单位为秒。	1-43200(12小时)范围内的某个值整数值，默认为30（秒）
	PollingWaitSeconds     int      `xml:"PollingWaitSeconds,omitempty"`     // 当 Queue 中没有消息时，针对该 Queue 的 ReceiveMessage 请求最长的等待时间，单位为秒。	0-30秒范围内的某个整数值，默认为0（秒）
	LoggingEnabled         string   `xml:"LoggingEnabled,omitempty"`         // 是否开启日志管理功能，True表示启用，False表示停用	True/False，默认为False
}

// QueueAttributesRequest queue attributes request
type QueueAttributesRequest struct {
	XMLName                xml.Name `xml:"Queue"`
	DelaySeconds           int      `xml:"DelaySeconds,omitempty"`           // 发送到该 Queue 的所有消息默认将以DelaySeconds参数指定的秒数延后可被消费，单位为秒。-- 0-604800秒（7天）范围内某个整数值，默认值为0
	MaximumMessageSize     int      `xml:"MaximumMessageSize,omitempty"`     // 发送到该Queue的消息体的最大长度，单位为byte。	1024(1KB)-65536（64KB）范围内的某个整数值，默认值为65536（64KB）。
	MessageRetentionPeriod int      `xml:"MessageRetentionPeriod,omitempty"` // 消息在该 Queue 中最长的存活时间，从发送到该队列开始经过此参数指定的时间后，不论消息是否被取出过都将被删除，单位为秒。	60 (1分钟)-1296000 (15 天)范围内某个整数值，默认值345600 (4 天)
	VisibilityTimeout      int      `xml:"VisibilityTimeout,omitempty"`      // 消息从该 Queue 中取出后从Active状态变成Inactive状态后的持续时间，单位为秒。	1-43200(12小时)范围内的某个值整数值，默认为30（秒）
	PollingWaitSeconds     int      `xml:"PollingWaitSeconds,omitempty"`     // 当 Queue 中没有消息时，针对该 Queue 的 ReceiveMessage 请求最长的等待时间，单位为秒。	0-30秒范围内的某个整数值，默认为0（秒）
	LoggingEnabled         string   `xml:"LoggingEnabled,omitempty"`         // 是否开启日志管理功能，True表示启用，False表示停用	True/False，默认为False
}

// SendMessageRequest send message request
type SendMessageRequest struct {
	XMLName      xml.Name `xml:"Message"`
	MessageBody  string   `xml:"MessageBody"`
	DelaySeconds int      `xml:"DelaySeconds,omitempty"`
	Priority     int      `xml:"Priority,omitempty"`
}

// BatchSendMessageRequest batch send message request
type BatchSendMessageRequest struct {
	XMLName  xml.Name             `xml:"Messages"`
	Messages []SendMessageRequest `xml:"Message"`
}

// BatchDeleteMessageRequest batch delete message request
type BatchDeleteMessageRequest struct {
	XMLName        xml.Name `xml:"Messages"`
	ReceiptHandles []string `xml:"ReceiptHandle"`
}

// CreateTopicRequest create topic request
type CreateTopicRequest struct {
	XMLName            xml.Name `xml:"Topic"`
	MaximumMessageSize int      `xml:"MaximumMessageSize,omitempty"`
	LoggingEnabled     string   `xml:"LoggingEnabled,omitempty"`
}

// TopicAttributesRequest topic attributes request
type TopicAttributesRequest struct {
	XMLName            xml.Name `xml:"Topic"`
	MaximumMessageSize int64    `xml:"MaximumMessageSize,omitempty"` //发送到该主题的消息体最大长度，单位为 Byte
	LoggingEnabled     string   `xml:"LoggingEnabled,omitempty"`     //是否开启日志管理功能，True表示启用，False表示停用
}

// CreateSubscribeRequest create subscribe request
type CreateSubscribeRequest struct {
	XMLName             xml.Name `xml:"Subscription"`
	EndPoint            string   `xml:"Endpoint"`                      // 描述此次订阅中接收消息的终端地址 -- 目前四种Endpoint: 1. HttpEndpoint，必须以”http://"为前缀 2. QueueEndpoint, 格式为acs:mns:{REGION}:{AccountID}:queues/{QueueName} 3. MailEndpoint, 格式为mail:directmail:{MailAddress} 4. SmsEndpoint, 格式为sms:directsms:anonymous 或sms:directsms:{Phone} -- Required
	FilterTag           string   `xml:"FilterTag,omitempty"`           // 描述了该订阅中消息过滤的标签（标签一致的消息才会被推送）-- 不超过16个字符的字符串，默认不进行消息过滤 -- Optional
	NotifyStrategy      string   `xml:"NotifyStrategy,omitempty"`      // 描述了向 Endpoint 推送消息出现错误时的重试策略 -- BACKOFF_RETRY 或者 EXPONENTIAL_DECAY_RETRY，默认为BACKOFF_RETRY，重试策略的具体描述请参考 基本概念/NotifyStrategy -- Optional
	NotifyContentFormat string   `xml:"NotifyContentFormat,omitempty"` // 描述了向 Endpoint 推送的消息格式 -- XML 、JSON 或者 SIMPLIFIED，默认为 XML，消息格式的具体描述请参考 基本概念/NotifyContentFormat -- Optional
}

// SubscribeAttributesRequest subscribe attributes request
type SubscribeAttributesRequest struct {
	XMLName        xml.Name `xml:"Subscription"`
	NotifyStrategy string   `xml:"NotifyStrategy,omitempty"` // 向 Endpoint 推送消息错误时的重试策略
}

// see https://help.aliyun.com/document_detail/27497.html?spm=5176.doc35141.6.722.ogncLw
type (
	// PublishMessageRequest publish message request
	PublishMessageRequest struct {
		XMLName           xml.Name          `xml:"Message"`
		MessageBody       string            `xml:"MessageBody"`
		MessageTag        string            `xml:"MessageTag,omitempty"`
		MessageAttributes MessageAttributes `xml:"MessageAttributes,omitempty"` // 消息属性，如果需要推送邮件或短信，则MessageAttributes为必填项
	}
	// DirectMail direct mail info
	DirectMail struct {
		Subject        string `json:"Subject"`
		AccountName    string `json:"AccountName"`
		ReplyToAddress int    `json:"ReplyToAddress"`
		AddressType    int    `json:"AddressType"`
		IsHTML         int    `json:"IsHtml"`
	}
	// DirectSMS direct sms info
	// see https://help.aliyun.com/document_detail/27497.html?spm=5176.doc51082.6.722.dOKA8e
	DirectSMS struct {
		FreeSignName string `json:"FreeSignName"`
		TemplateCode string `json:"TemplateCode"`
		Type         string `json:"Type"`
		Receiver     string `json:"Receiver"`
		SmsParams    string `json:"SmsParams"`
	}
	// MessageAttributes message attributes
	MessageAttributes struct {
		DirectMail string `xml:"DirectMail,omitempty"` // json
		DirectSMS  string `xml:"DirectSMS,omitempty"`  // json
	}
)
