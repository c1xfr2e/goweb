package mns

import "github.com/zaoshu/hermogo/proto"

// Client client interface
type Client interface {
	DoRequest(method string, path string, header map[string]string, req interface{}, resp interface{}) error
}

// Queue queue interface
type Queue interface {
	GetName() string
	Create(req proto.CreateQueueRequest) error
	IsExist() bool
	SetAttributes(req proto.QueueAttributesRequest) error
	Attributes() (proto.QueueAttributesResponse, error)
	Delete() error
	List(marker, number, prefix string) ([]string, string, error)
	SendMessage(req proto.SendMessageRequest) (proto.SendMessageResponse, error)
	BatchSendMessage(req []proto.SendMessageRequest) ([]proto.SendMessageResponse, error)
	ReceiveMessage(waitseconds int) (proto.ReceiveMessageResponse, error)
	BatchReceiveMessage(num, waitseconds int) ([]proto.ReceiveMessageResponse, error)
	PeekMessage() (proto.PeekMessageResponse, error)
	BatchPeekMessage(num int) ([]proto.PeekMessageResponse, error)
	DeleteMessage(receiptHandle string) error
	BatchDeleteMessage(receiptHandle []string) ([]proto.BatchDeleteMessageError, error)
	ChangeMessageVisibility(receiptHandle string, visibilityTimeout int) (proto.ChangeMessageVisibilityResponse, error)
}

// Topic topic interface
type Topic interface {
	GetName() string
	Create(req proto.CreateTopicRequest) error
	IsExist() bool
	SetAttributes(req proto.TopicAttributesRequest) error
	Attributes() (resp proto.TopicAttributesResponse, err error)
	Delete() error
	PublishMessage(req proto.PublishMessageRequest) (resp proto.PublishMessageResponse, err error)
	List(marker, number, prefix string) ([]string, string, error)
	Subscribe(subscribeName string, req proto.CreateSubscribeRequest) error
	SetSubscriptionAttributes(subscribeName string, req proto.SubscribeAttributesRequest) error
	SubscriptionAttributes(subscribeName string) (resp proto.SubscribeAttributesResponse, err error)
	Unsubscribe(subscribeName string) error
	ListSubscriptionByTopic(marker, number, prefix string) ([]string, string, error)
}
