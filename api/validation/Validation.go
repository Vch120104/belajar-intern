package validation

import (
	"after-sales/api/exceptions"
	"errors"
	"fmt"
	"net/http"

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
	// Inisialisasi validator dan translator sekali di fungsi init
	validate = validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, _ = uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(validate, trans)
}

func translateError(err error, trans ut.Translator) (errs []error) {
	if err == nil {
		return nil
	}
	validatorErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return []error{err}
	}
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(trans))
		errs = append(errs, translatedErr)
	}
	return errs
}

func ValidationForm(writer http.ResponseWriter, request *http.Request, form interface{}) *exceptions.BaseErrorResponse {
	err := validate.Struct(form)
	var msg string
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				msg = fmt.Sprintf("%s can't be empty", err.Field())
			case "email":
				msg = fmt.Sprintf("%s format not matched", err.Field())
			case "noWhiteSpace":
				msg = fmt.Sprintf("%s remove white space", err.Field())
			case "eqfield":
				msg = fmt.Sprintf("%s should match", err.Field())
			case "nefield":
				msg = fmt.Sprintf("%s shouldn't match", err.Field())
			default:
				msg = err.Translate(trans)
			}
		}
	}
	if msg != "" {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    msg,
			Err:        errors.New(msg),
		}
	} else if err != nil {
		errorMsg := fmt.Sprintf("%v ", translateError(err, trans))
		fmt.Println(err, " ++")
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    errorMsg,
			Err:        errors.New(errorMsg),
		}
	}

	return nil
}
