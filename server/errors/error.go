package errors

import (
	"net/http"
)

// Error list
var (
	ErrInternal              = New(http.StatusOK, -1, "Internal Server Error")
	ErrInvalidParameters     = New(http.StatusOK, 3, "Invalid Parameters")
	ErrUnauthorized          = New(http.StatusOK, 4, "Unauthorized")
	ErrInvalidNameOrPassword = New(http.StatusOK, 5, "Invalid Email or Password")
	ErrNoData          		 = New(http.StatusOK, 6, "No Data")
)
