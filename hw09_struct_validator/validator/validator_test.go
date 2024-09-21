package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	verr "github.com/AnnDutova/otus_go_hw/hw09_struct_validator/verrors"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:5"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	Email struct {
		Address string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	}

	UserWithEmail struct {
		ID    string `json:"id" validate:"len:36"`
		Email Email  `validate:"nested"`
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          nil,
			expectedErr: verr.ErrInvalidInputType,
		},
		{
			in: struct {
				Field bool `validate:"test"`
			}{},
			expectedErr: verr.ErrInvalidFieldType,
		},
		{
			in: struct {
				Field []bool `validate:"test"`
			}{},
			expectedErr: verr.ErrUnsupportedSliceType,
		},
		{
			in: struct {
				Field string `validate:"len::"`
			}{},
			expectedErr: verr.ErrInvalidRuleFormat,
		},
		{
			in: struct {
				Field string `validate:"length:4"`
			}{},
			expectedErr: verr.ErrUnsupportedAlias,
		},
		{
			in: struct {
				Field string `validate:"len:test"`
			}{},
			expectedErr: verr.ErrInvalidLength,
		},
		{
			in: struct {
				Field string `validate:"regexp:\\["`
			}{},
			expectedErr: verr.ErrInvalidRegexpValue,
		},
		{
			in: App{
				Version: "test",
			},
			expectedErr: verr.ErrUnexpectedLength,
		},
		{
			in:          Token{},
			expectedErr: nil,
		},
		{
			in:          Response{},
			expectedErr: verr.ErrUnmatchedIn,
		},
		{
			in: Response{
				Code: 200,
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 202,
			},
			expectedErr: verr.ErrUnmatchedIn,
		},
		{
			in: User{
				ID:     "123",
				Age:    123,
				Email:  "test",
				Role:   "test",
				Phones: []string{"123", "1234"},
			},
			expectedErr: verr.ValidationErrors{
				{
					Field: "ID",
					Value: "123",
					Err:   verr.ErrUnexpectedLength,
				},
				{
					Field: "Age",
					Value: 123,
					Err:   verr.ErrGreaterThanExpected,
				},
				{
					Field: "Email",
					Value: "test",
					Err:   verr.ErrUnmatchedRegexp,
				},
				{
					Field: "Role",
					Value: "test",
					Err:   verr.ErrUnmatchedIn,
				},
				{
					Field: "Phones",
					Value: "123",
					Err:   verr.ErrUnexpectedLength,
				},
				{
					Field: "Phones",
					Value: "1234",
					Err:   verr.ErrUnexpectedLength,
				},
			},
		},
		{
			in: User{
				ID:     "12345",
				Age:    23,
				Email:  "test@mail.com",
				Role:   "stuff",
				Phones: []string{"12345678910"},
			},
			expectedErr: nil,
		},
		{
			in: UserWithEmail{
				ID: "123",
				Email: Email{
					Address: "test",
				},
			},
			expectedErr: verr.ValidationErrors{
				{
					Field: "ID",
					Value: "123",
					Err:   verr.ErrUnexpectedLength,
				},
				{
					Field: "Address",
					Value: "test",
					Err:   verr.ErrUnmatchedRegexp,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			var vErr verr.ValidationErrors
			var expErr verr.ValidationErrors

			if errors.As(err, &vErr) {
				if errors.As(tt.expectedErr, &expErr) {
					assert.Equal(t, len(vErr), len(expErr))
					for i, verr := range vErr {
						assert.ErrorIs(t, verr.Err, expErr[i].Err)
					}
				}
			} else {
				assert.ErrorIs(t, err, tt.expectedErr)
			}

			if tt.expectedErr == nil {
				assert.NoError(t, err)
			}

			_ = tt
		})
	}
}
