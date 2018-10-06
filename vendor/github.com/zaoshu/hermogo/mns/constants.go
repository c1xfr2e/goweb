package mns

// APIVersion https://help.aliyun.com/document_detail/27485.html
const APIVersion = "2015-06-06"

// RequestContentType request header Content-Type
const RequestContentType = "text/xml;charset=utf-8"

// delay seconds
const (
	MinDelaySeconds     = 0
	MaxDelaySeconds     = 604800
	DefaultDelaySeconds = 0
)

// message size in byte
const (
	MinMessageSize     = 1024
	MaxMessageSize     = 65536 // 64k
	DefaultMessageSize = 65536
)

// message retention period in second
const (
	MinMessageRetentionPeriod     = 60
	MaxMessageRetentionPeriod     = 1296000 // 15 days
	DefaultMessageRetentionPeriod = 345600  // 4 days
)

// message visibility timeout in second
const (
	MinVisibilityTimeout     = 1
	MaxVisibilityTimeout     = 43200 // 12 hours
	DefaultVisibilityTimeout = 30
)

// polling wait seconds
const (
	MinPollingWaitSeconds     = 1
	MaxPollingWaitSeconds     = 30
	DefaultPollingWaitSeconds = 10
)

// batch receive message number
const (
	MinBatchReceiveMessageNumber     = 1
	MaxBatchReceiveMessageNumber     = 16
	DefaultBatchReceiveMessageNumber = 10
)

// batch peek message number
const (
	MinBatchPeekMessageNumber     = 1
	MaxBatchPeekMessageNumber     = 16
	DefaultBatchPeekMessageNumber = 10
)

// logging
const (
	LoggingEnabled  = "True"
	LoggingDisabled = "False"
)

// message priority
// see https://help.aliyun.com/document_detail/35134.html
const (
	MinMessagePriority     = 16
	MaxMessagePriority     = 1
	DefaultMessagePriority = 8
)

// create queue response http status code
// https://help.aliyun.com/document_detail/35129.html
const (
	CreateQueueSuccess            = 201
	CreateQueueSameStatusCode     = 204
	CreateQueueConfliceStatusCode = 409
)

// create topic subscribe notify strategy
const (
	BackoffRetryNotifyStrategy          = "BACKOFF_RETRY"
	ExponentialDecayRetryNotifyStrategy = "EXPONENTIAL_DECAY_RETRY"
)

// notify content format
const (
	XMLNotifyContentFormat        = "XML"
	JSONNotifyContentFormat       = "JSON"
	SIMPLIFIEDNotifyContentFormat = "SIMPLIFIED"
)
