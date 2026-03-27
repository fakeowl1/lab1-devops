package model

import "fmt"

type ApiError struct {
	Err error `json:"-"`
	
	Code int `json:"code"`
	Message string `json:"message"`
}

func NewApiError(err error, code int) (*ApiError) {
	return &ApiError{
		Err: err,
		Code: code,
		Message: err.Error(),
	}
} 

func (r *ApiError) Error() string {
	return fmt.Sprintf("%v", r.Err)
}
