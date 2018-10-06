package hermogo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zaoshu/hardcore/logging"
	"github.com/zaoshu/hardcore/trace"
)

type message struct {
	Meta map[string]string `json:"_meta"`
	Data json.RawMessage   `json:"_data"`
}

const requestIDKey = "X-Request-Id"

func encode(ctx context.Context, v interface{}) (context.Context, []byte, error) {
	var (
		data []byte
		err  error
	)

	if s, ok := v.(string); ok {
		// string value
		data = []byte(s)
	} else {
		data, err = json.Marshal(v)
		if err != nil {
			return ctx, nil, fmt.Errorf("json.Marshal failed, %v", err)
		}
	}

	requestID := trace.GetRequestID(ctx)
	if len(requestID) == 0 {
		requestID = trace.NewRequestID()
		ctx = logging.NewContext(ctx, logging.NewLoggerWithRequestID(requestID))
	}

	msg := message{
		Meta: map[string]string{
			requestIDKey: requestID,
		},
		Data: data,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return ctx, nil, fmt.Errorf("json.Marshal failed, %v", err)
	}
	return ctx, b, nil
}

func decode(v []byte) (ctx context.Context, data []byte, err error) {
	msg := message{}
	err = json.Unmarshal(v, &msg)
	if err != nil {
		// decode failed
		return
	}

	var requestID string
	if len(msg.Data) == 0 {
		// 旧的消息格式
		requestID = trace.NewRequestID()
		data = v
	} else {
		// 新的消息格式, 包含 _data, _meta
		if msg.Meta != nil {
			requestID = msg.Meta[requestIDKey]
		}
		if len(requestID) == 0 {
			requestID = trace.NewRequestID()
		}
		data = msg.Data
	}

	ctx = logging.NewContext(
		trace.NewContext(context.Background(), requestID),
		logging.NewLoggerWithRequestID(requestID),
	)
	return
}

func getRequestID(ctx context.Context) string {
	requestID := trace.GetRequestID(ctx)
	if len(requestID) == 0 {
		requestID = trace.NewRequestID()
	}
	return requestID
}
