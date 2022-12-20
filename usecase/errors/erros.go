package errors

import "errors"

var (
	ErrValidate   = errors.New("入力値が不正です")
	ErrPermission = errors.New("指定の操作を行う権限がありません")
)
