package proto

import (
	"encoding/xml"
	"fmt"
)

// see https://help.aliyun.com/document_detail/27501.html
const (
	ErrAccessDenied                = "AccessDenied"
	ErrInvalidAccessKeyID          = "InvalidAccessKeyId"
	ErrInternalError               = "InternalError"
	ErrInvalidAuthorizationHeader  = "InvalidAuthorizationHeader"
	ErrInvalidDateHeader           = "InvalidDateHeader"
	ErrInvalidArgument             = "InvalidArgument"
	ErrInvalidDegist               = "InvalidDegist"
	ErrInvalidRequestURL           = "InvalidRequestURL"
	ErrInvalidQueryString          = "InvalidQueryString"
	ErrMalformedXML                = "MalformedXML"
	ErrMissingAuthorizationHeader  = "MissingAuthorizationHeader"
	ErrMissingDateHeader           = "MissingDateHeader"
	ErrMissingReceiptHandle        = "MissingReceiptHandle"
	ErrMissingVisibilityTimeout    = "MissingVisibilityTimeout"
	ErrMessageNotExist             = "MessageNotExist"
	ErrQueueAlreadyExist           = "QueueAlreadyExist"
	ErrQueueDeletedRecently        = "QueueDeletedRecently"
	ErrInvalidQueueName            = "InvalidQueueName"
	ErrQueueNameLengthError        = "QueueNameLengthError"
	ErrQueueNotExist               = "QueueNotExist"
	ErrReceiptHandleError          = "ReceiptHandleError"
	ErrSignatureDoesNotMatch       = "SignatureDoesNotMatch"
	ErrTimeExpired                 = "TimeExpired"
	ErrQPSLimitExceeded            = "QpsLimitExceeded"
	ErrTopicAlreadyExist           = "TopicAlreadyExist"
	ErrTopicNameInvalid            = "TopicNameInvalid"
	ErrTopicNameLengthError        = "TopicNameLengthError"
	ErrTopicNotExist               = "TopicNotExist"
	ErrSubscriptionNameInvalid     = "SubscriptionNameInvalid"
	ErrSubscriptionNameLengthError = "SubscriptionNameLengthError"
	ErrSubscriptionNotExist        = "SubscriptionNotExist"
	ErrSubscriptionAlreadyExist    = "SubscriptionAlreadyExist"
	ErrEndpointInvalid             = "EndpointInvalid"
)

// Error mns result error
type Error struct {
	XMLName   xml.Name `xml:"Error"`
	Code      string   `xml:"Code"`
	Message   string   `xml:"Message"`
	RequestID string   `xml:"RequestID"`
	HostID    string   `xml:"HostID"`
}

// Error return string about error
func (e *Error) Error() string {
	return fmt.Sprintf(
		"Code: %s, Message: %s, RequestID: %s, HostID: %s",
		e.Code, e.Message, e.RequestID, e.HostID,
	)
}

// UnmarshalFromXML unmarshal from XML
func UnmarshalFromXML(data []byte, resp interface{}) error {
	err := xml.Unmarshal(data, resp)
	if err != nil {
		return fmt.Errorf(
			"xml.Unmarshal failed, reason %s, data %s",
			err, string(data),
		)
	}
	return nil
}

// NewErrorFromXML new error from xml
func NewErrorFromXML(data []byte) error {
	r := &Error{}
	err := UnmarshalFromXML(data, r)
	if err != nil {
		return err
	}
	return r
}

// IsError is error
func IsError(err error, code string) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*Error); ok && e.Code == code {
		return true
	}
	return false
}

// IsQueueNotExist is queue not exist
func IsQueueNotExist(err error) bool {
	return IsError(err, ErrQueueNotExist)
}

// IsTopicNotExist is topic not exist
func IsTopicNotExist(err error) bool {
	return IsError(err, ErrTopicNotExist)
}

// IsMessageNotExist is message not exist
func IsMessageNotExist(err error) bool {
	return IsError(err, ErrMessageNotExist)
}
