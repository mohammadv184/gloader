package data

import "errors"

var (
	ErrValueTypeParseFuncNotImplemented = errors.New("value type parse function not implemented")
	ErrValueTypeParentNotInitialized    = errors.New("value type parent not initialized")
	ErrDataTypeKindNotMatch             = errors.New("data type kind not match")
	ErrInvalidValue                     = errors.New("invalid passed value")
	ErrDestMustBePointer                = errors.New("destination must be a pointer")
	ErrDestNotAssignable                = errors.New("destination is not assignable")
)

var (
	ErrBufferAlreadyIsClosed = errors.New("buffer already is closed")
	ErrBufferIsClosed        = errors.New("buffer is closed")
)
