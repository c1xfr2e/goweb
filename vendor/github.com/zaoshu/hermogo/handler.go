package hermogo

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
)

type decodeError struct {
	Msg string
}

// Error error message
func (e decodeError) Error() string {
	return "decode error, " + e.Msg
}

// NewDecodeError new decode error
func NewDecodeError(msg string) error {
	return decodeError{msg}
}

// isDecodeError check if decode error
func isDecodeError(err error) bool {
	if err == nil {
		return false
	}
	if _, ok := err.(decodeError); ok {
		return true
	}
	return false
}

// MessageHandler message handler
type MessageHandler interface {
	// b: byte slice of message body
	Handle(ctx context.Context, b []byte) error
}

var _ MessageHandler = RawHandler{}

// RawHandlerFunc raw message handler function
type RawHandlerFunc func(ctx context.Context, b []byte) error

// RawHandler raw message handler
type RawHandler struct {
	handler RawHandlerFunc
}

// NewRawHandler new raw handler
func NewRawHandler(f RawHandlerFunc) MessageHandler {
	return RawHandler{handler: f}
}

// Handle call handler function
func (h RawHandler) Handle(ctx context.Context, b []byte) error {
	if b == nil {
		return nil
	}
	return h.handler(ctx, b)
}

// StructHandler struct handler, inject input struct
type StructHandler struct {
	method   reflect.Value
	argType  reflect.Type
	argIsPtr bool
}

// NewStructHandler handler: func(struct or *struct) error {}
func NewStructHandler(handler interface{}) (MessageHandler, error) {
	if handler == nil {
		return nil, errors.New("[NewStructHandler] handler must not be nil interface")
	}

	sh := StructHandler{
		method: reflect.ValueOf(handler),
	}

	if sh.method.Kind() != reflect.Func {
		return nil, errors.New("[NewStructHandler] handler Kind() must be reflect.Func")
	}
	if sh.method.IsNil() {
		return nil, errors.New("[NewStructHandler] value of handler interface must not be nil")
	}

	typ := reflect.TypeOf(handler)
	if typ.NumOut() != 1 {
		return nil, errors.New("[NewStructHandler] handler must have one error output parameter")
	}
	outType := typ.Out(0)
	if !(outType.Kind() == reflect.Interface &&
		outType.Implements(reflect.TypeOf((*error)(nil)).Elem())) {
		return nil, errors.New("[NewStructHandler] handler must have one error output parameter")
	}

	if typ.NumIn() != 2 {
		return nil, errors.New("[NewStructHandler] handler must have two input parameter, first parameter must be context, second parameter type must be struct or pointer to struct")
	}

	arg := typ.In(1)
	kind := arg.Kind()
	if kind == reflect.Ptr {
		kind = arg.Elem().Kind()
		sh.argIsPtr = true
		sh.argType = arg.Elem()
	} else {
		sh.argIsPtr = false
		sh.argType = arg
	}
	if kind != reflect.Struct {
		return nil, errors.New("[NewStructHandler] handler must have two input parameter, first parameter must be context, second parameter type must be struct or pointer to struct")
	}

	return sh, nil
}

// Handle call handler function
func (h StructHandler) Handle(ctx context.Context, b []byte) error {
	argv := reflect.New(h.argType) // argv is ptr
	if err := json.Unmarshal(b, argv.Interface()); err != nil {
		return NewDecodeError(err.Error())
	}

	ctxValue := reflect.ValueOf(ctx)
	var ret []reflect.Value
	if h.argIsPtr {
		ret = h.method.Call([]reflect.Value{ctxValue, argv})
	} else {
		ret = h.method.Call([]reflect.Value{ctxValue, argv.Elem()})
	}
	i := ret[0].Interface()
	if i == nil {
		return nil
	}
	return i.(error)
}
