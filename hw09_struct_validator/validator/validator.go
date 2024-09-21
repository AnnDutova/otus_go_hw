package validator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	verr "github.com/AnnDutova/otus_go_hw/hw09_struct_validator/verrors"
)

const (
	ValidateTag string = "validate"
	NestedTag   string = "nested"

	MinAlias    string = "min"
	MaxAlias    string = "max"
	InAlias     string = "in"
	LenAlias    string = "len"
	RegexpAlias string = "regexp"
)

func Validate(v interface{}) error {
	validationErrors := make(verr.ValidationErrors, 0)

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return verr.ErrInvalidInputType
	}

	if val.NumField() == 0 {
		return nil
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		alias, ok := field.Tag.Lookup(ValidateTag)
		if !ok {
			continue
		}

		if alias == NestedTag {
			err := Validate(val.Field(i).Interface())
			var valerr verr.ValidationErrors
			if errors.As(err, &valerr) {
				validationErrors = append(validationErrors, valerr...)
			} else {
				return err
			}
			continue
		}

		rules := strings.Split(alias, "|")

		for _, rule := range rules {
			err := executeValidation(field, val.Field(i), rule)

			var valerr verr.ValidationErrors
			if errors.As(err, &valerr) {
				validationErrors = append(validationErrors, valerr...)
			} else {
				return err
			}
		}
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func executeValidation(field reflect.StructField, value reflect.Value, rule string) error {
	validationErrors := make(verr.ValidationErrors, 0)

	switch field.Type.Kind() { //nolint
	case reflect.String:
		if err := executeStringRule(value, rule); err != nil {
			validationErrors = append(validationErrors,
				verr.ValidationError{
					Field: field.Name,
					Value: value,
					Err:   err,
				})
		}
	case reflect.Int:
		if err := executeIntRule(value, rule); err != nil {
			validationErrors = append(validationErrors,
				verr.ValidationError{
					Field: field.Name,
					Value: value,
					Err:   err,
				})
		}
	case reflect.Slice, reflect.Array:
		switch value.Interface().(type) {
		case []int:
			for i := 0; i < value.Len(); i++ {
				if err := executeIntRule(value.Index(i), rule); err != nil {
					validationErrors = append(validationErrors,
						verr.ValidationError{
							Field: field.Name,
							Value: value.Index(i),
							Err:   err,
						})
				}
			}
		case []string:
			for i := 0; i < value.Len(); i++ {
				if err := executeStringRule(value.Index(i), rule); err != nil {
					validationErrors = append(validationErrors,
						verr.ValidationError{
							Field: field.Name,
							Value: value.Index(i),
							Err:   err,
						})
				}
			}
		default:
			return verr.ErrUnsupportedSliceType
		}
	default:
		return verr.ErrInvalidFieldType
	}

	return validationErrors
}

func executeStringRule(value reflect.Value, rule string) error {
	aliasKey, aliasVal, err := splitAlias(rule)
	if err != nil {
		return err
	}

	switch aliasKey {
	case RegexpAlias, InAlias, LenAlias:
		return executeRule(value, aliasKey, aliasVal)
	default:
		return verr.ErrUnsupportedAlias
	}
}

func executeIntRule(value reflect.Value, rule string) error {
	aliasKey, aliasVal, err := splitAlias(rule)
	if err != nil {
		return err
	}

	switch aliasKey {
	case MinAlias, MaxAlias, InAlias:
		return executeRule(value, aliasKey, aliasVal)
	default:
		return verr.ErrUnsupportedAlias
	}
}

func executeRule(value reflect.Value, aliasKey, aliasVal string) error {
	switch aliasKey {
	case MinAlias, MaxAlias:
		return minMaxRule(value, aliasKey, aliasVal)
	case RegexpAlias:
		return regexpRule(value, aliasVal)
	case InAlias:
		return inRule(value, aliasVal)
	case LenAlias:
		return lenRule(value, aliasVal)
	}
	return nil
}

func splitAlias(rule string) (string, string, error) {
	ruleParts := strings.Split(rule, ":")
	if len(ruleParts) != 2 {
		return "", "", verr.ErrInvalidRuleFormat
	}

	return ruleParts[0], ruleParts[1], nil
}

func minMaxRule(value reflect.Value, aliasKey string, aliasVal string) error {
	alias, err := strconv.Atoi(aliasVal)
	if err != nil {
		return verr.ErrInvalidIntegerValue
	}
	if aliasKey == MinAlias {
		if int(value.Int()) < alias {
			return verr.ErrLessThanExpected
		}
		return nil
	}

	if int(value.Int()) > alias {
		return verr.ErrGreaterThanExpected
	}
	return nil
}

func regexpRule(value reflect.Value, aliasVal string) error {
	re, err := regexp.Compile(aliasVal)
	if err != nil {
		return verr.ErrInvalidRegexpValue
	}
	if !re.MatchString(value.String()) {
		return verr.ErrUnmatchedRegexp
	}
	return nil
}

func lenRule(value reflect.Value, aliasVal string) error {
	length, err := strconv.Atoi(aliasVal)
	if err != nil {
		return verr.ErrInvalidLength
	}
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		if value.Len() != length {
			return verr.ErrUnexpectedLength
		}
	} else if value.Kind() == reflect.String {
		if len(value.String()) != length {
			return verr.ErrUnexpectedLength
		}
	}
	return nil
}

func inRule(value reflect.Value, aliasVal string) error {
	found := false
	options := strings.Split(aliasVal, ",")

	if value.Kind() == reflect.String {
		for _, option := range options {
			if value.String() == option {
				found = true
				break
			}
		}
	} else if value.Kind() == reflect.Int {
		for _, option := range options {
			val, err := strconv.Atoi(option)
			if err != nil {
				return verr.ErrInvalidLength
			}

			if int(value.Int()) == val {
				found = true
				break
			}
		}
	}

	if !found {
		return verr.ErrUnmatchedIn
	}
	return nil
}
