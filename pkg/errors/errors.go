package er

import (
	"fmt"
)

var (
	// rest api
	ErrorExist    = New(1, "запись существует")
	ErrorNotExist = New(2, "запись не существует")
)

type Error struct {
	ErrorCode        int    `json:"error_code"`
	ErrorDescription string `json:"error_description"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("err_type:%d , err_des:%s", e.Code(), e.ErrorDescription)
}
func (e *Error) Code() int {
	return e.ErrorCode
}

func New(code int, desc string) *Error {
	return &Error{
		ErrorCode:        code,
		ErrorDescription: desc,
	}
}
