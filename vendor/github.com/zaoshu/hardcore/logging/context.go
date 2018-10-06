package logging

import (
	"context"

	"github.com/sirupsen/logrus"
)

// ContextLoggerKey key of logger
const ContextLoggerKey = "__zaoshu/hardcore/logging/key/logger__"

// FromContext new logger from context
func FromContext(ctx context.Context) *logrus.Entry {
	logger, ok := ctx.Value(ContextLoggerKey).(*logrus.Entry)
	if ok {
		return logger
	}
	return logrus.NewEntry(logrus.StandardLogger())
}

// NewContext new context with logger
func NewContext(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, ContextLoggerKey, logger)
}
