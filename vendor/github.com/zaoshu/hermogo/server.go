package hermogo

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sync"
	"time"

	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/zaoshu/hardcore/logging"
	"github.com/zaoshu/hermogo/mns"
	"github.com/zaoshu/hermogo/proto"
)

type queueHandler struct {
	queue   string
	handler MessageHandler

	maxConcurrency int
	concurrency    int
	cond           *sync.Cond
	wg             *sync.WaitGroup

	mu       sync.Mutex
	needStop bool
}

// start start handler
func (h *queueHandler) start() error {
	logrus.Infof("[hermogo] start to listen queue %s, waitseconds %d", h.queue, WaitSeconds)

	q, err := mns.NewMNSQueue(h.queue, defaultClient)
	if err != nil {
		return err
	}
	go func() {
		h.loop(q)
	}()
	return nil
}

// Stop stop handler
func (h *queueHandler) stop() {
	h.mu.Lock()
	if h.needStop {
		h.mu.Unlock()
		return
	}
	h.needStop = true
	h.mu.Unlock()

	h.wg.Wait()
	logrus.Infof("[hermogo] server of queue %s stopped", h.queue)
}

func (h *queueHandler) shouldStop() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.needStop
}

// receiveMessage receive message
func (h *queueHandler) loop(q mns.Queue) {
	for !h.shouldStop() {
		var num int
		h.cond.L.Lock()
		if h.concurrency >= h.maxConcurrency {
			h.cond.Wait()
		}
		num = h.maxConcurrency - h.concurrency
		h.cond.L.Unlock()

		if num < 1 || h.shouldStop() {
			continue
		} else if num > mns.MaxBatchReceiveMessageNumber {
			num = mns.MaxBatchReceiveMessageNumber
		}

		messages, err := q.BatchReceiveMessage(num, WaitSeconds)
		if err != nil {
			if proto.IsMessageNotExist(err) {
				// no message in queue, continue
			} else if proto.IsQueueNotExist(err) {
				err = q.Create(mns.NewCreateQueueRequest())
				if err != nil && !proto.IsError(err, proto.ErrQueueAlreadyExist) {
					logrus.Errorf("[hermogo] create queue %s failed in receive message loop, %v", h.queue, err)
					randomSleepSeconds()
				}
			} else {
				logrus.Errorf("[hermogo] receive message failed in loop, queue %s, %v", h.queue, err)
				randomSleepSeconds()
			}
			continue
		}

		// handle messages
		h.cond.L.Lock()
		h.concurrency += len(messages)
		h.cond.L.Unlock()
		h.wg.Add(len(messages))
		for _, msg := range messages {
			go h.handleMessage(q, msg)
		}
	}

	logrus.Infof("[hermogo] message loop of queue %s stopped", h.queue)
}

func (h *queueHandler) handleMessage(q mns.Queue, msg proto.ReceiveMessageResponse) {
	defer func() {
		h.cond.L.Lock()
		h.concurrency -= 1
		h.cond.L.Unlock()
		h.cond.Signal()
		h.wg.Done()
	}()

	ctx, body, err := decode([]byte(msg.MessageBody))
	if err != nil {
		logrus.Errorf("[hermogo] decode message `%s` failed and will delete this message, queue %s, %v", msg.MessageBody, q.GetName(), err)
		// decode failed, delete this message
		err = q.DeleteMessage(msg.ReceiptHandle)
		if err != nil {
			logrus.Errorf("[hermogo] delete message failed after decode failed, queue %s, %v", q.GetName(), err)
		}
		return
	}

	logger := logging.FromContext(ctx)
	logger.Infof("[hermogo] received message `%s` from queue %s", string(body), h.queue)

	start := time.Now()
	defer func() {
		latency := uint(time.Since(start) / time.Millisecond)
		logger.WithFields(logging.NewStatsField(
			"mq",
			map[string]interface{}{
				"name":    q.GetName(),
				"action":  "process",
				"latency": latency,
			},
		)).Infof("[hermogo] message handled of queue %s, latency %d", h.queue, latency)
	}()

	err = handle(ctx, h.handler, body)
	if err != nil {
		if isDecodeError(err) {
			// decode error
			logger.Errorf("[hermogo] decode failed when handle message, delete this message, %v", err)
			err = q.DeleteMessage(msg.ReceiptHandle)
			if err != nil {
				logger.Errorf("[hermogo] delete message failed after handler decode failed, %v", err)
			}
		} else {
			logger.Errorf("[hermogo] handle message failed, do not delete it, %v", err)
		}
	} else {
		// no error, need to delete this message
		err = q.DeleteMessage(msg.ReceiptHandle)
		if err != nil {
			logger.Errorf("[hermogo] handle message success but delete failed, %v", err)
		}
	}
	return
}

// call handler
func handle(ctx context.Context, handler MessageHandler, b []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch p := r.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("panic when call message handler, %v", p)
			}
		}
	}()

	err = handler.Handle(ctx, b)
	return
}

var serverNamePattern *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]{0,64}$`)

// Server message queue server
type Server struct {
	name     string // server name, also topic subscriber name
	handlers map[string]*queueHandler
	mu       sync.Mutex
	wg       sync.WaitGroup
}

// NewServer new message queue server
//
// name: server name, also topic subscriber name, pattern is ^[a-zA-Z][a-zA-Z0-9-]{0,64}$
func NewServer(name string) (*Server, error) {
	if defaultClient == nil {
		return nil, errors.New("default client not initialized")
	}
	if !serverNamePattern.MatchString(name) {
		return nil, errors.New("invalid server name: " + name)
	}
	return &Server{
		name:     name,
		handlers: map[string]*queueHandler{},
	}, nil
}

// AddHandlers add struct handlers
//
// handlers: slice of Handler
func (s *Server) AddHandlers(handlers ...Handler) error {
	for _, h := range handlers {
		err := s.AddStructHandler(h.Queue, h.HandlerFunc, h.Concurrency)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddStructHandler add struct handler
//
// queue: queue name
// h: handler function
// concurrency: max concurrency, default is 1, max is 1000
func (s *Server) AddStructHandler(queue string, h interface{}, concurrency ...int) error {
	sh, err := NewStructHandler(h)
	if err != nil {
		return err
	}
	return s.addHandler(queue, sh, concurrency...)
}

// AddRawHandler add raw handler
//
// queue: queue name
// h: handler function
// concurrency: max concurrency, default is 1, max is 1000
func (s *Server) AddRawHandler(queue string, h RawHandlerFunc, concurrency ...int) error {
	return s.addHandler(queue, NewRawHandler(h), concurrency...)
}

func (s *Server) addHandler(queue string, h MessageHandler, concurrency ...int) error {
	var con int
	if len(concurrency) > 0 {
		con = concurrency[0]
	}
	if con < 1 {
		con = 1
	} else if con > 1000 {
		con = 1000
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[queue] = &queueHandler{
		queue:          queue,
		handler:        h,
		maxConcurrency: con,
		cond:           sync.NewCond(&sync.Mutex{}),
		wg:             &sync.WaitGroup{},
	}
	return nil
}

// Subscribe subscribe
//
// subs: slice of Subscriber
func (s *Server) Subscribe(subs ...Subscriber) error {
	for _, sub := range subs {
		err := s.subscribe(sub.Topic, sub.HandlerFunc, sub.Concurrency)
		if err != nil {
			return err
		}
	}
	return nil
}

// Subscribe subscribe topic
//
// topic:       topic name
// h: struct handler function
// concurrency: max concurrency, default is 1, max is 1000
func (s *Server) subscribe(topic string, h interface{}, concurrency ...int) error {
	if len(topic) == 0 {
		return errors.New("topic name cannot be empty")
	}

	queue := fmt.Sprintf("%s-%s", topic, s.name)
	endPoint, err := s.endPoint(queue)
	if err != nil {
		return fmt.Errorf("get end point queue of topic %s failed, %v", topic, err)
	}

	// create queue
	q, err := mns.NewMNSQueue(queue, defaultClient)
	if err != nil {
		return err
	}
	err = q.Create(mns.NewCreateQueueRequest())
	if err != nil && !proto.IsError(err, proto.ErrQueueAlreadyExist) {
		return err
	}

	// create topic
	t, err := mns.NewMNSTopic(topic, defaultClient)
	if err != nil {
		return err
	}
	err = t.Create(mns.NewCreateTopicRequest())
	if err != nil && !proto.IsError(err, proto.ErrTopicAlreadyExist) {
		return err
	}

	// subscribe
	err = t.Subscribe(s.name, mns.NewCreateSubscribeRequest(endPoint, ""))
	if err != nil && !proto.IsError(err, proto.ErrSubscriptionAlreadyExist) {
		return err
	}

	return s.AddStructHandler(queue, h, concurrency...)
}

// Run run server
func (s *Server) Run() error {
	if len(s.handlers) == 0 {
		return nil
	}

	s.wg.Add(len(s.handlers))
	for _, h := range s.handlers {
		err := h.start()
		if err != nil {
			return err
		}
	}
	go func() {
		ex := make(chan os.Signal, 1)
		signal.Notify(ex, syscall.SIGTERM, syscall.SIGINT)
		<-ex
		logrus.Info("[hermogo] server shutting down...")
		s.Stop()
	}()
	s.wg.Wait()
	logrus.Info("[hermogo] server stopped")
	return nil
}

// Stop stop server
func (s *Server) Stop() {
	for _, h := range s.handlers {
		h.stop()
		s.wg.Done()
	}
}

func (s *Server) endPoint(queue string) (string, error) {
	return fmt.Sprintf("acs:mns:%s:%s:queues/%s", defaultConfig.region, defaultConfig.accountID, queue), nil
}
