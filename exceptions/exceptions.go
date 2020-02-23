package exceptions

import "errors"

var (
	//ErrDBHandle DB构建失败
	ErrDBHandle = errors.New("handle db error")
)
