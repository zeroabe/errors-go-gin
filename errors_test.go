package ginerrors

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
)

type testCase struct {
	input    interface{}
	hasError bool
}

func TestValidationErrors(t *testing.T) {
	type diveStruct struct {
		Email string `json:"email_val" validate:"email"`
		Min   string `json:"min_val_test" validate:"min=1"`
		Max   string `json:"max_val_test" validate:"max=4"`
	}
	type testStruct struct {
		DiveTest []diveStruct `json:"dive_test" validate:"min=1,dive"`
	}

	type testCase struct {
		name  string
		value testStruct
	}
	var (
		v     = validator.New()
		cases = []testCase{
			{
				name:  "zero dive",
				value: testStruct{},
			},
			{
				name: "one dive with invalid email",
				value: testStruct{
					DiveTest: []diveStruct{
						{Email: "invalid", Min: "xxx", Max: "xxx"},
					},
				},
			},
			{
				name: "one dive with invalid email and max",
				value: testStruct{
					DiveTest: []diveStruct{
						{Email: "invalid", Min: "xxx", Max: "xxxxx"},
					},
				},
			},
		}
	)

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			verr := v.Struct(cc.value)

			_, obj := makeResponse(verr, "en")
			assert.NotNil(t, obj.Validation)
		})
	}
}

func TestHasErrors(t *testing.T) {
	tcs := map[string]testCase{
		"common error": {input: errors.New("some error"), hasError: true},
		"errors slice": {input: []error{errors.New("some error"), errors.New("some more error")}, hasError: true},
		"errors map":   {input: map[string]error{"some error": errors.New("some error"), "some more error": errors.New("some more error")}, hasError: true},
		//"gorm error":   {input: gorm.Errors{errors.New("some error"), errors.New("some more error")}, hasError: true},
	}

	for caseName, tc := range tcs {
		if HasErrors(tc.input) != tc.hasError {
			fmt.Printf("test fail in `%s`", caseName)
			t.FailNow()
		}
	}
}

func TestNew(t *testing.T) {
	errMsg := "err msg"
	err := New(errMsg)
	if err == nil {
		fmt.Print("no errors occurs")
		t.FailNow()
	}
	if errMsg != err.Error() {
		fmt.Print("error message is not equals to err.Error() string")
		t.FailNow()
	}
}

func TestNewf(t *testing.T) {
	errFormat := "err msg: %s | %s"
	errMsg1 := "alert 1!"
	errMsg2 := "alert 2!"
	err := Newf(errFormat, errMsg1, errMsg2)
	if err == nil {
		fmt.Print("no errors occurs")
		t.FailNow()
	}

	if fmt.Sprintf(errFormat, errMsg1, errMsg2) != err.Error() {
		fmt.Print("errors are not equals")
		t.FailNow()
	}
}
