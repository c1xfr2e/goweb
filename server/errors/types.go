package errors

import "fmt"

// Error represents business error
// 	HTTPStatus for http code
// 	Code for business error code
//  Message for error msg
type Error struct {
	HTTPStatus int    `json:"httpStatus"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
}

// Error implements error interface
func (err Error) Error() string {
	return fmt.Sprintf("[%d][%d] %s", err.HTTPStatus, err.Code, err.Message)
}

func New(status, code int, msg string) error {
	return &Error{
		HTTPStatus: status,
		Code:       code,
		Message:    msg,
	}
}
