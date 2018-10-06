package proto

import "encoding/xml"

// QueueAttributesResponse queue attributes
type QueueAttributesResponse struct {
	XMLName                xml.Name `xml:"Queue"`                            //
	QueueName              string   `xml:"QueueName,omitempty"`              //Queue 的名称
	CreateTime             int64    `xml:"CreateTime,omitempty"`             //Queue 的创建时间，从1970-1-1 00:00:00 到现在的秒值
	LastModifyTime         int64    `xml:"LastModifyTime,omitempty"`         //修改 Queue 属性信息最近时间，从1970-1-1 00:00:00 到现在的秒值
	DelaySeconds           int      `xml:"DelaySeconds,omitempty"`           //发送消息到该 Queue 的所有消息默认将以 DelaySeconds 参数指定的秒数延后可被消费，单位为秒
	MaximumMessageSize     int      `xml:"MaximumMessageSize,omitempty"`     //发送到该 Queue 的消息体的最大长度，单位为byte
	MessageRetentionPeriod int      `xml:"MessageRetentionPeriod,omitempty"` //消息在该 Queue 中最长的存活时间，从发送到该队列开始经过此参数指定的时间后，不论消息是否被取出过都将被删除，单位为秒
	PollingWaitSeconds     int      `xml:"PollingWaitSeconds,omitempty"`     //当 Queue 消息量为空时，针对该 Queue 的 ReceiveMessage 请求最长的等待时间，单位为秒
	ActiveMessages         int      `xml:"Activemessages,omitempty"`         //在该 Queue 中处于 Active 状态的消息总数，为近似值
	InactiveMessages       int      `xml:"InactiveMessages,omitempty"`       //在该 Queue 中处于 Inactive 状态的消息总数，为近似值
	DelayMessages          int      `xml:"DelayMessages,omitempty"`          //在该 Queue 中处于 Delayed 状态的消息总数，为近似值
	LoggingEnabled         string   `xml:"LoggingEnabled,omitempty"`         // 是否开启日志管理功能，True表示启用，False表示停用
}

// TopicAttributesResponse topic attributes response
type TopicAttributesResponse struct {
	XMLName                xml.Name `xml:"Topic"`
	TopicName              string   `xml:"TopicName"`              //主题名称
	CreateTime             int64    `xml:"CreateTime"`             //主题的创建时间，从 1970-1-1 00:00:00到现在的秒值
	LastModifyTime         int64    `xml:"LastModifyTime"`         //修改主题属性信息的最近时间，从 1970-1-1 00:00:00 到现在的秒值
	MaximumMessageSize     int64    `xml:"MaximumMessageSize"`     //发送到该主题的消息体最大长度，单位为 Byte
	MessageRetentionPeriod int64    `xml:"MessageRetentionPeriod"` //消息在主题中最长存活时间，从发送到该主题开始经过此参数指定的时间后，不论消息是否被成功推送给用户都将被删除，单位为秒
	MessageCount           int64    `xml:"MessageCount"`           //当前该主题中消息数目
	LoggingEnabled         bool     `xml:"LoggingEnabled"`         //是否开启日志管理功能，True表示启用，False表示停用
}

// SubscribeAttributesResponse subscribe attributes response
type SubscribeAttributesResponse struct {
	XMLName             xml.Name `xml:"Subscription"`
	SubscriptionName    string   // Subscription 的名称
	Subscriber          string   // Subscription 订阅者的 AccountId
	TopicOwner          string   // Subscription 订阅的主题所有者的 AccountId
	TopicName           string   // Subscription 订阅的主题名称
	Endpoint            string   // 订阅的终端地址
	NotifyStrategy      string   // 向 Endpoint 推送消息错误时的重试策略
	NotifyContentFormat string   // 向 Endpoint 推送的消息内容格式
	FilterTag           string   // 描述了该订阅中消息过滤的标签（仅标签一致的消息才会被推送）
	CreateTime          int64    // Subscription 的创建时间，从 1970-1-1 00:00:00 到现在的秒值
	LastModifyTime      int64    // 修改 Subscription 属性信息最近时间，从 1970-1-1 00:00:00 到现在的秒值
}

type (
	// QueueInfo queue info
	QueueInfo struct {
		XMLName xml.Name `xml:"Queue"`
		URL     string   `xml:"QueueURL"`
	}
	// QueueListResponse queue list response
	QueueListResponse struct {
		XMLName    xml.Name    `xml:"Queues"`
		List       []QueueInfo `xml:"Queue"`
		NextMarker string      `xml:"NextMarker"`
	}
)

type (
	// TopicInfo topic info
	TopicInfo struct {
		XMLName xml.Name `xml:"Topic"`
		URL     string   `xml:"TopicURL"`
	}
	// TopicListResponse topic list response
	TopicListResponse struct {
		XMLName    xml.Name    `xml:"Topics"`
		List       []TopicInfo `xml:"Topic"`
		NextMarker string      `xml:"NextMarker"`
	}
)

type (
	// SubscriptionInfo subscription info
	SubscriptionInfo struct {
		XMLName xml.Name `xml:"Subscription"`
		URL     string   `xml:"SubscriptionURL"`
	}
	// SubscriptionListResponse subscription list response
	SubscriptionListResponse struct {
		XMLName    xml.Name           `xml:"Subscriptions"`
		List       []SubscriptionInfo `xml:"Subscription"`
		NextMarker string             `xml:"NextMarker"`
	}
)

// SendMessageResponse send message response
type SendMessageResponse struct {
	XMLName        xml.Name `xml:"Message"`
	MessageID      string   `xml:"MessageId"`
	MessageBodyMD5 string   `xml:"MessageBodyMD5"`
	ReceiptHandle  string   `xml:"ReceiptHandle"`
}

// BatchSendMessageResponse send message response
type BatchSendMessageResponse struct {
	XMLName  xml.Name              `xml:"Messages"`
	Messages []SendMessageResponse `xml:"Message"`
}

// ReceiveMessageResponse receive message response
type ReceiveMessageResponse struct {
	XMLName          xml.Name `xml:"Message"`
	MessageID        string   `xml:"MessageId"`
	ReceiptHandle    string   `xml:"ReceiptHandle"`
	MessageBody      string   `xml:"MessageBody"`
	MessageBodyMD5   string   `xml:"MessageBodyMD5"`
	EnqueueTime      int64    `xml:"EnqueueTime"`
	NextVisibleTime  int64    `xml:"NextVisibleTime"`
	FirstDequeueTime int64    `xml:"FirstDequeueTime"`
	DequeueCount     int      `xml:"DequeueCount"`
	Priority         int      `xml:"Priority"`
}

// BatchReceiveMessageResponse batch receive message response
type BatchReceiveMessageResponse struct {
	XMLName  xml.Name                 `xml:"Messages"`
	Messages []ReceiveMessageResponse `xml:"Message"`
}

// ChangeMessageVisibilityResponse response
type ChangeMessageVisibilityResponse struct {
	XMLName         xml.Name `xml:"ChangeVisibility"`
	ReceiptHandle   string   `xml:"ReceiptHandle"`
	NextVisibleTime int64    `xml:"NextVisibleTime"`
}

type (
	// BatchDeleteMessageError batch delete message error
	BatchDeleteMessageError struct {
		XMLName       xml.Name `xml:"Error"`
		ErrorCode     string   `xml:"ErrorCode"`
		ErrorMessage  string   `xml:"ErrorMessage"`
		ReceiptHandle string   `xml:"ReceiptHandle"`
	}
	// BatchDeleteMessageResponse batch delete message response
	BatchDeleteMessageResponse struct {
		XMLName xml.Name                  `xml:"Errors"`
		Errors  []BatchDeleteMessageError `xml:"Error"`
	}
)

// PeekMessageResponse peek message response
type PeekMessageResponse struct {
	XMLName          xml.Name `xml:"Message"`
	MessageID        string   `xml:"MessageId"`
	MessageBodyMD5   string   `xml:"MessageBodyMD5"`
	MessageBody      string   `xml:"MessageBody"`
	EnqueueTime      int64    `xml:"EnqueueTime"`
	FirstDequeueTime int64    `xml:"FirstDequeueTime"`
	DequeueCount     int      `xml:"DequeueCount"`
	Priority         int      `xml:"Priority"`
}

// BatchPeekMessageResponse batch peek message response
type BatchPeekMessageResponse struct {
	XMLName  xml.Name              `xml:"Messages"`
	Messages []PeekMessageResponse `xml:"Message"`
}

// PublishMessageResponse publish message response
type PublishMessageResponse struct {
	XMLName        xml.Name `xml:"Message"`
	MessageID      string   `xml:"MessageId"`
	MessageBodyMD5 string   `xml:"MessageBodyMD5"`
}
