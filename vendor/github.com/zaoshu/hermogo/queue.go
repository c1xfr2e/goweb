package hermogo

import (
	"context"
	"errors"

	"time"

	"github.com/zaoshu/hardcore/logging"
	"github.com/zaoshu/hermogo/mns"
	"github.com/zaoshu/hermogo/proto"
)

// Send send queue message
//
// queue: name of queue
// v:     string or any thing that can be encoded to json
// delay: seconds of delay, max 604800
//
// if queue not exist, create new queue, then send message
func Send(ctx context.Context, queue string, v interface{}, delay ...int) error {
	if defaultClient == nil {
		return errors.New("default client not initialized")
	}

	q, err := mns.NewMNSQueue(queue, defaultClient)
	if err != nil {
		return err
	}
	return SendWithQueue(ctx, q, v, delay...)
}

// QueueAttributes get queue attribute
// queue: name of queue
// return:
// 	resp: queue attribute
func QueueAttributes(ctx context.Context, queue string) (resp proto.QueueAttributesResponse, err error) {
	if defaultClient == nil {
		err = errors.New("default client not initialized")
		return
	}

	q, err := mns.NewMNSQueue(queue, defaultClient)
	if err != nil {
		return
	}
	return q.Attributes()
}

// SendWithQueue send with queue
func SendWithQueue(ctx context.Context, q mns.Queue, v interface{}, delay ...int) error {
	ctx, body, err := encode(ctx, v)
	if err != nil {
		return err
	}

	var delaySeconds int
	if len(delay) > 0 {
		delaySeconds = delay[0]
	}

	for i := 0; i < 3; i++ {
		start := time.Now()
		resp, err := q.SendMessage(mns.NewSendMessageRequest(string(body), delaySeconds))
		if err != nil {
			if proto.IsQueueNotExist(err) {
				// 队列不存在
				err = q.Create(mns.NewCreateQueueRequest())
				if err == nil || proto.IsError(err, proto.ErrQueueAlreadyExist) {
					continue
				}
			}
			logging.FromContext(ctx).Errorf("[hermogo] send message `%s` failed, queue %s, %v", body, q.GetName(), err)
			return err
		}
		logging.FromContext(ctx).WithFields(logging.NewStatsField(
			"mq",
			map[string]interface{}{
				"name":    q.GetName(),
				"action":  "send",
				"latency": uint(time.Since(start) / time.Millisecond),
			},
		)).Infof("[hermogo] send message `%s`, queue %s, messageID %s", body, q.GetName(), resp.MessageID)
		return nil
	}
	return errors.New("[hermogo] send message failed, should not be here")
}
