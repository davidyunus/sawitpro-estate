package domain

import "errors"

var (
	ErrInvalidInput   = errors.New("invalid input")
	ErrMaxSizeEstate  = errors.New("max size of state exceed 50000")
	ErrLocationFilled = errors.New("location already filled")
	ErrEstateNotFound = errors.New("estate not found")
)
