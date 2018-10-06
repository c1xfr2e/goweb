package hermogo

import (
	"math/rand"
	"time"

	"github.com/zaoshu/hermogo/proto"
)

// CanRetry check if can retry
func CanRetry(err error) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*proto.Error); ok {
		// mns 服务器返回的错误
		// 只有服务器内部错误可以重试
		return e.Code == proto.ErrInternalError
	}
	// 其他错误都可以重试
	return true
}

func randomSleepSeconds() {
	time.Sleep(time.Second * time.Duration(rand.Intn(5)+3))
}
