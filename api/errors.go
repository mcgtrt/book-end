package api

import "net/http"

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e Error) Error() string {
	return e.Msg
}

func NewError(code int, msg string) Error {
	return Error{
		Code: code,
		Msg:  msg,
	}
}

func ErrInvalidCredentials() Error {
	return Error{
		Code: http.StatusBadRequest,
		Msg:  "invalid credentials",
	}
}

func ErrResourceNotFound() Error {
	return Error{
		Code: http.StatusBadRequest,
		Msg:  "resource not found",
	}
}

func ErrInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Msg:  "invalid id",
	}
}

func ErrUnauthorised() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Msg:  "unauthorised",
	}
}

func ErrBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Msg:  "bad request",
	}
}

func ErrInternalServerError() Error {
	return Error{
		Code: http.StatusInternalServerError,
		Msg:  "internal server error",
	}
}
