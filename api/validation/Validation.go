package validation

import (
	"after-sales/api/exceptions"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	validate = validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, _ = uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(validate, trans)

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func ValidationForm(writer http.ResponseWriter, request *http.Request, form interface{}) *exceptions.BaseErrorResponse {
	err := validate.Struct(form)
	if err == nil {
		return nil
	}

	var msgBuilder strings.Builder
	for _, err := range err.(validator.ValidationErrors) {
		field := strings.Replace(err.Field(), "_", " ", -1)
		msgBuilder.WriteString(getErrorMessage(err.Tag(), field))
	}

	msg := msgBuilder.String()
	return &exceptions.BaseErrorResponse{
		StatusCode: http.StatusBadRequest,
		Message:    msg,
		Err:        errors.New(msg),
	}
}

func getErrorMessage(tag string, fieldName string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s can't be empty; ", fieldName)
	case "email":
		return fmt.Sprintf("%s format is invalid; ", fieldName)
	case "noWhiteSpace":
		return fmt.Sprintf("%s should not contain whitespace; ", fieldName)
	case "eqfield":
		return fmt.Sprintf("%s should match; ", fieldName)
	case "nefield":
		return fmt.Sprintf("%s should not match; ", fieldName)
	default:
		return fmt.Sprintf("%s is invalid; ", fieldName)
	}
}
