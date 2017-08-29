package mchain

import (
	"errors"
)

func RecoverIntoError(pointerToError *error) {
	if r := recover(); r != nil {
		switch x := r.(type) {
		case error:
			*pointerToError = x
		case string:
			*pointerToError = errors.New(x)
		default:
			*pointerToError = errors.New("unknown panic")
		}
	}
}
