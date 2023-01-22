package data

import "errors"

var (
	ErrParseFuncNotImplemented = errors.New("parse function not implemented")
	ErrDataTypeKindNotMatch    = errors.New("data type kind not match")
	ErrInvalidValue            = errors.New("invalid passed value")
	ErrDestMustBePointer       = errors.New("destination must be a pointer")
	ErrDestNotAssignable       = errors.New("destination is not assignable")
)

var (
	ErrBufferAlreadyIsClosed = errors.New("buffer already is closed")
	ErrBufferIsClosed        = errors.New("buffer is closed")
)
