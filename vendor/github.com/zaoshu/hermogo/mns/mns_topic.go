package mns

import (
	"errors"

	"github.com/zaoshu/hermogo/proto"
)

// MNSTopic mns topic
type MNSTopic struct {
	name   string
	client Client
}

var _ Topic = MNSTopic{}

// NewMNSTopic new mns topic
func NewMNSTopic(name string, client Client) (Topic, error) {
	if !mnsQueueNamePattern.MatchString(name) {
		return nil, errors.New("invalid topic name " + name)
	}
	if client == nil {
		return nil, errors.New("client must not be nil")
	}
	return MNSTopic{name, client}, nil
}

// GetName get topic name
func (t MNSTopic) GetName() string {
	return t.name
}

// Create create mns topic, maybe return error QueueAlreadyExist
func (t MNSTopic) Create(req proto.CreateTopicRequest) error {
	return t.client.DoRequest("PUT", "/topics/"+t.name, nil, req, nil)
}

// IsExist check queue is exist
func (t MNSTopic) IsExist() bool {
	_, err := t.Attributes()
	if proto.IsTopicNotExist(err) {
		return false
	} else if err != nil {
		return false
	} else {
		return true
	}
}

// SetAttributes set topic attributes
func (t MNSTopic) SetAttributes(req proto.TopicAttributesRequest) error {
	return t.client.DoRequest("PUT", "/topics/"+t.name+"?metaoverride=true", nil, req, nil)
}

// Attributes get topic attributes
func (t MNSTopic) Attributes() (resp proto.TopicAttributesResponse, err error) {
	err = t.client.DoRequest("GET", "/topics/"+t.name, nil, nil, &resp)
	return
}

// Delete delete topic
func (t MNSTopic) Delete() error {
	return t.client.DoRequest("DELETE", "/topics/"+t.name, nil, nil, nil)
}

// PublishMessage pulish message by topic
func (t MNSTopic) PublishMessage(req proto.PublishMessageRequest) (resp proto.PublishMessageResponse, err error) {
	err = t.client.DoRequest("POST", "/topics/"+t.name+"/messages", nil, req, &resp)
	return
}

// List get topic list by marker, number and prefix
func (t MNSTopic) List(marker, number, prefix string) ([]string, string, error) {
	header := map[string]string{}
	if len(marker) > 0 {
		header["x-mns-marker"] = marker
	}
	if len(number) > 0 {
		header["x-mns-ret-number"] = number
	}
	if len(prefix) > 0 {
		header["x-mns-prefix"] = prefix
	}

	resp := proto.TopicListResponse{}
	err := t.client.DoRequest("GET", "/topics", header, nil, &resp)
	if err != nil {
		return nil, "", err
	}

	result := []string{}
	for _, value := range resp.List {
		result = append(result, value.URL)
	}
	return result, resp.NextMarker, nil
}

// Subscribe subscribe to topic
func (t MNSTopic) Subscribe(subscribeName string, req proto.CreateSubscribeRequest) error {
	if len(req.FilterTag) > 16 {
		return errors.New("The length of filter tag should be between 1 and 16")
	}
	return t.client.DoRequest("PUT", "/topics/"+t.name+"/subscriptions/"+subscribeName, nil, req, nil)
}

// SetSubscriptionAttributes set subscription attributes
func (t MNSTopic) SetSubscriptionAttributes(subscribeName string, req proto.SubscribeAttributesRequest) error {
	return t.client.DoRequest("PUT", "/topics/"+t.name+"/subscriptions/"+subscribeName+"?metaoverride=true", nil, req, nil)
}

// SubscriptionAttributes get subscription attributes
func (t MNSTopic) SubscriptionAttributes(subscribeName string) (resp proto.SubscribeAttributesResponse, err error) {
	err = t.client.DoRequest("GET", "/topics/"+t.name+"/subscriptions/"+subscribeName, nil, nil, &resp)
	return
}

// Unsubscribe unsubscribe topic
func (t MNSTopic) Unsubscribe(subscribeName string) error {
	return t.client.DoRequest("DELETE", "/topics/"+t.name+"/subscriptions/"+subscribeName, nil, nil, nil)
}

// ListSubscriptionByTopic get subscription list by topic
func (t MNSTopic) ListSubscriptionByTopic(marker, number, prefix string) ([]string, string, error) {
	header := map[string]string{}
	if len(marker) > 0 {
		header["x-mns-marker"] = marker
	}
	if len(number) > 0 {
		header["x-mns-ret-number"] = number
	}
	if len(prefix) > 0 {
		header["x-mns-prefix"] = prefix
	}

	resp := proto.SubscriptionListResponse{}
	err := t.client.DoRequest("GET", "/topics/"+t.name+"/subscriptions", header, nil, &resp)
	if err != nil {
		return nil, "", err
	}

	result := []string{}
	for _, value := range resp.List {
		result = append(result, value.URL)
	}
	return result, resp.NextMarker, nil
}

// NewCreateTopicRequest new create topic request
func NewCreateTopicRequest() proto.CreateTopicRequest {
	return proto.CreateTopicRequest{
		MaximumMessageSize: DefaultMessageSize,
		LoggingEnabled:     LoggingEnabled,
	}
}

// NewCreateSubscribeRequest new create subscribe request
func NewCreateSubscribeRequest(endpoint, filterTag string) proto.CreateSubscribeRequest {
	return proto.CreateSubscribeRequest{
		EndPoint:            endpoint,
		FilterTag:           filterTag,
		NotifyStrategy:      ExponentialDecayRetryNotifyStrategy,
		NotifyContentFormat: SIMPLIFIEDNotifyContentFormat,
	}
}

// NewPublishMessageRequest new publish message request
func NewPublishMessageRequest(body, tag string) proto.PublishMessageRequest {
	return proto.PublishMessageRequest{
		MessageBody: body,
		MessageTag:  tag,
	}
}
