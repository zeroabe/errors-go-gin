package ginerrors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"gopkg.in/go-playground/validator.v9"
)

const namespaceSeparator = "."

type langName string
type validationRule string
type errorPattern string
type validationErrors map[validationRule]errorPattern

func (ve errorPattern) string() string {
	return string(ve)
}

var (
	CommonValidationErrors = map[langName]validationErrors{
		"ru": {
			"ek":       "Ошибка валидации для свойства `%s` с правилом `%s`",
			"required": "Свойство `%s` обязательно для заполнения",
			"gt":       "Свойство `%s` должно содержать более `%s` элементов",
		},
		"en": {
			"ek":       "field `%s` is invalid: `%s`",
			"required": "field `%s` is required",
			"gt":       "field `%s` must contain more than `%s` elements",
			"email":    "field `%s` is not valid email",
			"min":      "length of `%s` field value is shorter then `%s`",
			"max":      "length of `%s` field value is greater then `%s`",
		},
	}
)

var (
	defaultLang = "en"

	ErrNotFound       = errors.New("route not found")
	ErrNoMethod       = errors.New("method not allowed")
	ErrServerError    = errors.New("internal server error")
	ErrRecordNotFound = errors.New("record not found")
)

func getLang(c *gin.Context) langName {
	lang := c.GetHeader("lang")
	if lang == "" {
		lang = c.DefaultQuery("lang", defaultLang)
	}

	return langName(lang)
}

// validationErrors Формирование массива ошибок
func makeErrorsSlice(err interface{}, lang langName, level int) map[FieldName]interface{} {
	ve := make(map[FieldName]interface{})
	for _, e := range err.(validator.ValidationErrors) {
		fieldName, vee := processFieldError(e, lang)
		keys := splitNamespace(e.Namespace())[1:]

		cur := mapWalk(ve, keys)
		if cur == nil {
			continue
		}

		cur[fieldName] = vee
	}

	return ve
}

func mapWalk(m map[FieldName]interface{}, keys []FieldName) map[FieldName]interface{} {
	var (
		ok  bool
		cur map[FieldName]interface{}
	)

	for i, k := range keys {
		if i == 0 {
			if _, ok := m[k]; !ok {
				m[k] = make(map[FieldName]interface{})
			}

			if cur, ok = m[k].(map[FieldName]interface{}); !ok {
				return nil
			}

			continue
		}

		if _, ok := cur[k]; !ok {
			if cur != nil {
				cur[k] = make(map[FieldName]interface{})
				if cur, ok = cur[k].(map[FieldName]interface{}); !ok {
					return nil
				}
			}

			continue
		}

		if cur, ok = cur[k].(map[FieldName]interface{}); !ok {
			return nil
		}
	}

	return cur
}

func processFieldError(e validator.FieldError, lang langName) (FieldName, interface{}) {
	field := getFieldName(e.Namespace())
	er := getErrMessage(validationRule(e.ActualTag()), field, e.Param(), lang)

	return field, er
}

func makeErrorsSliceFromViolations(violations []*errdetails.BadRequest_FieldViolation) map[FieldName]interface{} {
	ve := make(map[FieldName]interface{})
	for _, v := range violations {
		if v == nil {
			continue
		}
		field := FieldName(v.Field)
		if _, ok := ve[field]; !ok {
			ve[field] = make([]ValidationError, 0)
		}
	}

	return ve
}

func splitNamespace(ns string) []FieldName {
	ns = strings.Replace(ns, "]", "", -1)
	ns = strings.Replace(ns, "[", namespaceSeparator, -1)
	values := strings.Split(ns, namespaceSeparator)

	result := make([]FieldName, 0)
	for _, k := range values {
		result = append(result, FieldName(k))
	}

	return result
}

func getFieldName(namespace string) FieldName {
	namespaceSlice := splitNamespace(namespace)
	fieldName := namespaceSlice[len(namespaceSlice)-1]

	return fieldName
}

func getErrMessage(errorType validationRule, field FieldName, param string, lang langName) ValidationError {
	errKey := errorType
	_, ok := CommonValidationErrors[lang][errorType]
	if !ok {
		errKey = "ek"
	}

	if param != "" && errKey == "ek" {
		return ValidationError(fmt.Sprintf(CommonValidationErrors[lang][errKey].string(), field, param))
	}

	params := []interface{}{field}
	if param != "" {
		params = append(params, param)
	}

	return ValidationError(fmt.Sprintf(CommonValidationErrors[lang][errKey].string(), params...))
}
