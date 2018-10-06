package hermogo

import (
	"context"
	"errors"

	"github.com/zaoshu/hardcore/logging"
	"github.com/zaoshu/hermogo/mns"
	"github.com/zaoshu/hermogo/proto"
)

// PublishMessage send topic message
//
// topic: name of topic
// v:     string or any thing that can be encoded to json
//
// if topic not exists, create new topic, then publish message
func PublishMessage(ctx context.Context, topic string, v interface{}) error {
	if defaultClient == nil {
		return errors.New("default client not initialized")
	}

	t, err := mns.NewMNSTopic(topic, defaultClient)
	if err != nil {
		return err
	}

	ctx, body, err := encode(ctx, v)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		resp, err := t.PublishMessage(mns.NewPublishMessageRequest(string(body), ""))
		if err != nil {
			if proto.IsTopicNotExist(err) {
				err = t.Create(mns.NewCreateTopicRequest())
				if err == nil || proto.IsError(err, proto.ErrTopicAlreadyExist) {
					continue
				}
			}
			logging.FromContext(ctx).Errorf("[hermogo] publish message `%s` failed, topic %s", body, t.GetName())
			return err
		}
		logging.FromContext(ctx).Infof("[hermogo] publish message `%s`, topic %s, messageID %s", body, t.GetName(), resp.MessageID)
		return nil
	}
	return errors.New("[hermogo] publish message failed, should not be here")
}
