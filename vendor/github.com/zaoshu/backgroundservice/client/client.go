package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zaoshu/backgroundservice/types/constants"
	"github.com/zaoshu/backgroundservice/types/request"
	"github.com/zaoshu/hermogo"
)

func sendTask(ctx context.Context, handler string, argument interface{}) error {
	b, err := json.Marshal(argument)
	if err != nil {
		return fmt.Errorf("json.Marshal argument failed, %v", err)
	}
	msg := request.TaskRequest{
		Handler:       handler,
		Expire:        -1,
		ArgumentsJSON: string(b),
	}
	return hermogo.Send(ctx, constants.BackgroundServiceQueue, msg)
}

// SendEmail send new email task
func SendEmail(ctx context.Context, subject, body string, receivers []string, replies []string, bucket string, attachments map[string]string) error {
	argument := request.MailArgument{
		Body:        body,
		Subject:     subject,
		Receivers:   receivers,
		Bucket:      bucket,
		ReplyTo:     replies,
		Attachments: attachments,
	}
	return sendTask(ctx, constants.MailHandler, argument)
}

// NewSignupEmailTask create sending signup email task
func NewSignupEmailTask(ctx context.Context, signupURL, receiver string) error {
	argument := request.MailSignupEmailArgument{URL: signupURL, Receiver: receiver}
	return sendTask(ctx, constants.MailSignupHandler, argument)
}

// NewResetPasswordEmailTask create sending reset password email task
func NewResetPasswordEmailTask(ctx context.Context, resetURL, receiver string) error {
	argument := request.MailResetPasswordEmailArgument{URL: resetURL, Receiver: receiver}
	return sendTask(ctx, constants.MailResetPasswordHandler, argument)
}

// NewTaskNoticeEmailTask create task notice email task
func NewTaskNoticeEmailTask(ctx context.Context, status, finishTimeString, taskName, taskStatusURL string, receivers []string) error {
	argument := request.MailTaskNoticeEmailArgument{
		Status:     status,
		FinishTime: finishTimeString,
		TaskName:   taskName,
		TaskURL:    taskStatusURL,
		Receivers:  receivers,
	}
	return sendTask(ctx, constants.MailTaskNoticeHandler, argument)
}
