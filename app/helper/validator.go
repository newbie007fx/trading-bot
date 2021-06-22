package helper

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func getValidate() *validator.Validate {
	if validate == nil {
		validate = validator.New()
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			return name
		})
	}

	return validate
}

func Validate(request interface{}) (error, map[string]interface{}) {
	validate := getValidate()

	err := validate.Struct(request)
	if err == nil {
		return nil, structToMap(request)
	}

	return NewValidationError(err), map[string]interface{}{}
}

func structToMap(data interface{}) map[string]interface{} {
	m := make(map[string]interface{})

	elem := reflect.ValueOf(data).Elem()
	relType := elem.Type()

	for i := 0; i < relType.NumField(); i++ {
		m[relType.Field(i).Tag.Get("json")] = elem.Field(i).Interface()
	}

	return m
}

type ValidationError struct {
	err error
	msg map[string]string
}

func NewValidationError(err error) *ValidationError {
	ver := &ValidationError{err: err}
	ver.TranslateMsg()

	return ver
}

func (ver *ValidationError) TranslateMsg() {
	errMsg := map[string]string{}

	if v, ok := ver.err.(validator.ValidationErrors); ok {
		for _, fieldError := range v {
			errMsg[fieldError.Field()] = GetValidationMsg(fieldError.Field(), fieldError.Tag())
		}
	}

	ver.msg = errMsg
}

func (ver *ValidationError) GetErrorMsg() map[string]string {
	return ver.msg
}

func (ValidationError) Error() string {
	return "validation error"
}
