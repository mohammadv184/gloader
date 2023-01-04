package data

import "errors"

var (
	ErrParseFuncNotImplemented = errors.New("parse function not implemented")
	ErrDataTypeKindNotMatch    = errors.New("data type kind not match")
	ErrInvalidValue            = errors.New("invalid passed value")
)
