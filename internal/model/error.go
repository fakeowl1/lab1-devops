package model

import "fmt"

type ApiError struct {
	Err error `json:"-"`
	
	Code string `json:"code"`
	Message string `json:"message"`
}

func NewApiError(err error, code string) (*ApiError) {
	return &ApiError{
		Err: err,
		Code: code,
		Message: err.Error(),
	}
} 

func (r *ApiError) Error() string {
	return fmt.Sprintf("%v", r.Err)
}
