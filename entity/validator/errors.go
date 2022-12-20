package validator

import "errors"

var (
	ErrInvalidType   = errors.New("型が相違しています")
	ErrRequired      = errors.New("必須フィールドが欠落しています")
	ErrSpecialRune   = errors.New("特殊文字が含まれています")
	ErrInvalidRune   = errors.New("無効な文字が含まれています")
	ErrInvalidFormat = errors.New("無効な形式です")
	ErrInvalidLength = errors.New("長さが指定の範囲外です")
	ErrForbidden     = errors.New("使用できないリソースを指定しています")
)
