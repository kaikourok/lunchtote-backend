package validator

import (
	"reflect"
)

/*-------------------------------------------------------------------------------------------------
	IsNotOnlySpace
-------------------------------------------------------------------------------------------------*/

var IsNotOnlySpace = NotOnlySpaceRule{}

type NotOnlySpaceRule struct{}

func (r NotOnlySpaceRule) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.String {
		return ErrInvalidType
	}

	str := v.String()
	if str != "" && isOnlySpaces(str) {
		return ErrRequired
	}

	return nil
}

/*-------------------------------------------------------------------------------------------------
	IsNotContainSpecialRunes
-------------------------------------------------------------------------------------------------*/

var IsNotContainSpecialRune = NotContainSpecialRuneRule{}

type NotContainSpecialRuneRule struct{}

func (r NotContainSpecialRuneRule) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.String {
		return ErrInvalidType
	}

	str := v.String()
	if isContainSpecialRunes(str) {
		return ErrSpecialRune
	}

	return nil
}
