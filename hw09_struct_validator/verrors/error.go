package verrors

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidInputType     = errors.New("invalid input type")
	ErrInvalidFieldType     = errors.New("invalid field type")
	ErrUnsupportedSliceType = errors.New("unsupported slice/array type")
	ErrInvalidRuleFormat    = errors.New("invalid rule format")

	ErrUnsupportedAlias = errors.New("unsupported validation alias")

	ErrInvalidLength    = errors.New("invalid length value in validator field")
	ErrUnexpectedLength = errors.New("value length does not match expected length")

	ErrInvalidIntegerValue = errors.New("invalid integer value in validator field")
	ErrLessThanExpected    = errors.New("value is less than expected minimum")
	ErrGreaterThanExpected = errors.New("value is greater than expected maximum")

	ErrInvalidRegexpValue = errors.New("invalid regexp in validator field")
	ErrUnmatchedRegexp    = errors.New("value does not match regexp rule")

	ErrUnmatchedIn = errors.New("value does not match in rule")
)

type ValidationError struct {
	Field string
	Value any
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	buf := strings.Builder{}
	for _, err := range v {
		buf.WriteString(fmt.Sprintf("%s = %v: %s\n", err.Field, err.Value, err.Err))
	}
	return buf.String()
}
