package xerr

import (
	"github.com/zeromicro/x/errors"
)

/**
@Author: loser
@Description: to create new error
*/

// @Description: this interface(CodeMsg) implment the error
func New(code int, msg string) error {
	return errors.New(code, msg)
}

func NewMsg(msg string) error {
	return errors.New(SERVER_COMMON_ERROR, msg)
}

func NewDBErr() error {
	return errors.New(DB_ERROR, ErrMsg(DB_ERROR))
}

func NewInternalErr() error {
	return errors.New(SERVER_COMMON_ERROR, ErrMsg(SERVER_COMMON_ERROR))
}
