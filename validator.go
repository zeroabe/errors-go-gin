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
			"ek":       "Fail to validate field `%s` with rule `%s`",
			"required": "field `%s` is required",
			"gt":       "field `%s` must contain more than `%s` elements",
			"email":    "field `%s` is not valid email",
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
		field, val := processFieldError(e, lang, level)
		if v, ok := val.(map[FieldName]interface{}); ok {
			ve[field] = v[field]
			continue
		}

		ve[field] = val
	}

	return ve
}

func processFieldError(e validator.FieldError, lang langName, level int) (FieldName, interface{}) {
	field := getFieldName(e.Namespace(), level)
	isNested, nextLevel := isNested(e.Namespace(), level)

	if !isNested || nextLevel == 0 {
		er := make([]ValidationError, 0)
		er = append(
			er,
			getErrMessage(validationRule(e.ActualTag()), field, e.ActualTag(), lang),
		)

		return field, er
	}

	ve := make(map[FieldName]interface{})
	if _, ok := ve[field]; !ok {
		ve[field] = make(map[FieldName]interface{})
	}

	nextField, value := processFieldError(e, lang, nextLevel)
	vve := make(map[FieldName]interface{})
	vve[nextField] = value
	ve[field] = vve

	return field, ve
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
		//e := ValidationError(v.Description)
		//ve[field] = append(ve[field], e)
	}

	return ve
}

func isNested(namespace string, level int) (bool, int) {
	parts := strings.Split(namespace, namespaceSeparator)
	return len(parts) > (2 - level), (len(parts) - 1) - (level + 1)
}

func getFieldName(namespace string, level int) FieldName {
	namespace = strings.Replace(namespace, "]", "", -1)
	namespace = strings.Replace(namespace, "[", ".", -1)
	namespaceSlice := strings.Split(namespace, ".")
	fieldName := strings.ToLower(namespaceSlice[level+1])

	//if len(namespaceSlice) > 2 {
	//	fieldName = strings.Join([]string{strings.Join(namespaceSlice[1:len(namespaceSlice)-1], "."), field}, ".")
	//}

	return FieldName(fieldName)
}

func getErrMessage(errorType validationRule, field FieldName, param string, lang langName) ValidationError {
	errKey := errorType
	_, ok := CommonValidationErrors[lang][errorType]
	if !ok {
		errKey = "ek"
	}

	if param != "" && errKey == "ek" {
		return ValidationError(fmt.Sprintf(CommonValidationErrors[lang][errKey].string(), field, errorType))
	}

	return ValidationError(fmt.Sprintf(CommonValidationErrors[lang][errKey].string(), field))
}
