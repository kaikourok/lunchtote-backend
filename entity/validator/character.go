package validator

import (
	"reflect"

	"github.com/kaikourok/lunchtote-backend/entity/service"
)

/*-------------------------------------------------------------------------------------------------
	IsUsername
-------------------------------------------------------------------------------------------------*/

type UsernameRule struct {
	minLength int
	maxLength int
}

func (r UsernameRule) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.String {
		return ErrInvalidType
	}

	str := v.String()
	if len(str) < r.minLength || r.maxLength < len(str) {
		return ErrInvalidLength
	}

	for _, r := range str {
		if !(r == '_') && !('0' <= r && r <= '9') && !('A' <= r && r <= 'Z') && !('a' <= r && r <= 'z') {
			return ErrInvalidRune
		}
	}

	return nil
}

func IsUsername(minLength, maxLength int) UsernameRule {
	if minLength < 1 {
		panic("usernameの最小長は1以上でなければいけません")
	}
	if maxLength < minLength {
		panic("usernameのmaxLengthよりminLengthが長くなっています")
	}
	return UsernameRule{minLength, maxLength}
}

/*-------------------------------------------------------------------------------------------------
	Password
-------------------------------------------------------------------------------------------------*/

var IsPassword = PasswordRule{}

type PasswordRule struct{}

func (r PasswordRule) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.String {
		return ErrInvalidType
	}

	str := v.String()
	if str == "" {
		return ErrRequired
	}

	if isHexadecimal(str) {
		return nil
	} else {
		return ErrInvalidRune
	}
}

/*-------------------------------------------------------------------------------------------------
	ImagePath
-------------------------------------------------------------------------------------------------*/

type ImagePathRule struct {
	character *int
}

func (r ImagePathRule) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.String {
		return ErrInvalidType
	}

	str := v.String()
	characterId, err := service.ParseFilePath(str)
	if err != nil {
		return ErrInvalidFormat
	}

	if r.character != nil && characterId != *r.character {
		return ErrForbidden
	}

	return nil
}

func IsImagePath(characterId *int) ImagePathRule {
	return ImagePathRule{characterId}
}

/*-------------------------------------------------------------------------------------------------
	ImagePathOrEmpty
-------------------------------------------------------------------------------------------------*/

type ImagePathOrEmptyRule struct {
	character *int
}

func (r ImagePathOrEmptyRule) Validate(value interface{}) error {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.String {
		return ErrInvalidType
	}

	str := v.String()
	if str == "" {
		return nil
	}

	characterId, err := service.ParseFilePath(str)
	if err != nil {
		return ErrInvalidFormat
	}

	if r.character != nil && characterId != *r.character {
		return ErrForbidden
	}

	return nil
}

func IsImagePathOrEmpty(characterId *int) ImagePathOrEmptyRule {
	return ImagePathOrEmptyRule{characterId}
}
