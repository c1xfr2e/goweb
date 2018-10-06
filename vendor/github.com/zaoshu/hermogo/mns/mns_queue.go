package mns

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/zaoshu/hermogo/proto"
)

var mnsQueueNamePattern *regexp.Regexp

func init() {
	mnsQueueNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]{0,255}$`)
}

// MNSQueue mns queue
type MNSQueue struct {
	name   string
	client Client
}

var _ Queue = MNSQueue{}

// NewMNSQueue new mns queue
func NewMNSQueue(name string, client Client) (Queue, error) {
	if !mnsQueueNamePattern.MatchString(name) {
		return nil, errors.New("invalid queue name " + name)
	}
	if client == nil {
		return nil, errors.New("client must not be nil")
	}
	return MNSQueue{name, client}, nil
}

// GetName get queue name
func (q MNSQueue) GetName() string {
	return q.name
}

// Create create mns queue, maybe return error QueueAlreadyExist
func (q MNSQueue) Create(req proto.CreateQueueRequest) error {
	return q.client.DoRequest("PUT", "/queues/"+q.name, nil, req, nil)
}

// IsExist check queue is exist
func (q MNSQueue) IsExist() bool {
	_, err := q.Attributes()
	if proto.IsQueueNotExist(err) {
		return false
	} else if err != nil {
		return false
	} else {
		return true
	}
}

// SetAttributes set queue attributes
func (q MNSQueue) SetAttributes(req proto.QueueAttributesRequest) error {
	return q.client.DoRequest("PUT", "/queues/"+q.name+"?metaoverride=true", nil, req, nil)
}

// Attributes get queue attributes
func (q MNSQueue) Attributes() (resp proto.QueueAttributesResponse, err error) {
	err = q.client.DoRequest("GET", "/queues/"+q.name, nil, nil, &resp)
	return
}

// Delete delete queue
func (q MNSQueue) Delete() error {
	return q.client.DoRequest("DELETE", "/queues/"+q.name, nil, nil, nil)
}

// List get queue list
func (q MNSQueue) List(marker, number, prefix string) ([]string, string, error) {
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

	resp := proto.QueueListResponse{}
	err := q.client.DoRequest("GET", "/queues", header, nil, &resp)
	if err != nil {
		return nil, "", err
	}

	result := []string{}
	for _, value := range resp.List {
		result = append(result, value.URL)
	}
	return result, resp.NextMarker, nil
}

// SendMessage send message
func (q MNSQueue) SendMessage(req proto.SendMessageRequest) (resp proto.SendMessageResponse, err error) {
	err = q.client.DoRequest("POST", "/queues/"+q.name+"/messages", nil, req, &resp)
	return
}

// BatchSendMessage batch send message
func (q MNSQueue) BatchSendMessage(req []proto.SendMessageRequest) ([]proto.SendMessageResponse, error) {
	reqs := proto.BatchSendMessageRequest{
		Messages: req,
	}
	resp := proto.BatchSendMessageResponse{}
	err := q.client.DoRequest("POST", "/queues/"+q.name+"/messages", nil, reqs, &resp)
	return resp.Messages, err
}

// ReceiveMessage receive message
func (q MNSQueue) ReceiveMessage(waitseconds int) (resp proto.ReceiveMessageResponse, err error) {
	if waitseconds < MinPollingWaitSeconds {
		waitseconds = DefaultPollingWaitSeconds
	} else if waitseconds > MaxPollingWaitSeconds {
		waitseconds = MaxPollingWaitSeconds
	}
	err = q.client.DoRequest(
		"GET",
		fmt.Sprintf("/queues/%s/messages?waitseconds=%d", q.name, waitseconds),
		nil, nil, &resp,
	)
	return
}

// BatchReceiveMessage batch receive message
func (q MNSQueue) BatchReceiveMessage(num, waitseconds int) ([]proto.ReceiveMessageResponse, error) {
	if num < MinBatchReceiveMessageNumber {
		num = DefaultBatchReceiveMessageNumber
	} else if num > MaxBatchReceiveMessageNumber {
		num = MaxBatchReceiveMessageNumber
	}

	if waitseconds < MinPollingWaitSeconds {
		waitseconds = DefaultPollingWaitSeconds
	} else if waitseconds > MaxPollingWaitSeconds {
		waitseconds = MaxPollingWaitSeconds
	}

	resp := proto.BatchReceiveMessageResponse{}
	path := fmt.Sprintf("/queues/%s/messages?numOfMessages=%d&waitseconds=%d", q.name, num, waitseconds)
	err := q.client.DoRequest("GET", path, nil, nil, &resp)
	return resp.Messages, err
}

// DeleteMessage delete message
func (q MNSQueue) DeleteMessage(receiptHandle string) error {
	return q.client.DoRequest(
		"DELETE",
		fmt.Sprintf("/queues/%s/messages?ReceiptHandle=%s", q.name, receiptHandle),
		nil, nil, nil,
	)
}

// BatchDeleteMessage batch delete message
func (q MNSQueue) BatchDeleteMessage(receiptHandle []string) ([]proto.BatchDeleteMessageError, error) {
	req := proto.BatchDeleteMessageRequest{
		ReceiptHandles: receiptHandle,
	}
	resp := proto.BatchDeleteMessageResponse{}
	err := q.client.DoRequest("DELETE", "/queues/"+q.name+"/messages", nil, req, &resp)
	return resp.Errors, err
}

// PeekMessage peek message
func (q MNSQueue) PeekMessage() (resp proto.PeekMessageResponse, err error) {
	err = q.client.DoRequest("GET", "/queues/"+q.name+"/messages?peekonly=true", nil, nil, &resp)
	return
}

// BatchPeekMessage batch peek message
func (q MNSQueue) BatchPeekMessage(num int) ([]proto.PeekMessageResponse, error) {
	if num < MinBatchPeekMessageNumber {
		num = DefaultBatchPeekMessageNumber
	} else if num > MaxBatchPeekMessageNumber {
		num = MaxBatchPeekMessageNumber
	}
	resp := proto.BatchPeekMessageResponse{}
	err := q.client.DoRequest("GET", fmt.Sprintf("/queues/%s/messages?peekonly=true&numOfMessages=%d", q.name, num), nil, nil, &resp)
	return resp.Messages, err
}

// ChangeMessageVisibility change message visibility time
func (q MNSQueue) ChangeMessageVisibility(receiptHandle string, visibilityTimeout int) (resp proto.ChangeMessageVisibilityResponse, err error) {
	path := fmt.Sprintf("/queues/%s/messages?ReceiptHandle=%s&visibilityTimeout=%d", q.name, receiptHandle, visibilityTimeout)
	err = q.client.DoRequest("PUT", path, nil, nil, &resp)
	return
}

// NewCreateQueueRequest default create queue request
func NewCreateQueueRequest() proto.CreateQueueRequest {
	return proto.CreateQueueRequest{
		MessageRetentionPeriod: MaxMessageRetentionPeriod,
		VisibilityTimeout:      2 * 60, // 2 min
		LoggingEnabled:         LoggingEnabled,
	}
}

// NewSendMessageRequest default send message request
func NewSendMessageRequest(body string, delaySeconds int) proto.SendMessageRequest {
	return proto.SendMessageRequest{
		MessageBody:  body,
		DelaySeconds: delaySeconds,
		Priority:     DefaultMessagePriority,
	}
}
