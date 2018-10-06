package hermogo

// Handler handler of queue message
type Handler struct {
	Queue       string
	HandlerFunc interface{} // function like func(ctx context.Context, req TestStruct) error
	Concurrency int         // default is 1, max 1000
}

// Subscriber subscribe
type Subscriber struct {
	Topic       string
	HandlerFunc interface{} // function like func(ctx context.Context, req TestStruct) error
	Concurrency int         // default is 1, max 1000
}
