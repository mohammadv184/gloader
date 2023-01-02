package types

import (
	"fmt"
	"gloader/data"
)

type StringType struct {
	data.BaseType
}

func (t *StringType) parseType(p []byte) error {
	fmt.Println("StringType.parseType() called")
	return nil
}
func (t *StringType) To(tt data.Type) data.Type {
	fmt.Println("StringType.To() called")
	return nil
}
func (t *StringType) GetValue() []byte {
	fmt.Println("StringType.GetValue() called")
	return nil
}
