package validation

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func validateUsername(fl validator.FieldLevel) bool {
	return usernameRegex.MatchString(fl.Field().String())
}

func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("username", validateUsername)
	}
}

func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "username":
		return fmt.Sprintf("The %s field can only contain [a-z, A-Z, 0-9, - and _]", fe.Field())
	case "required":
		return fmt.Sprintf("The %s field cannot be blank", fe.Field())
	case "min":
		return fmt.Sprintf("The %s field min length is %s", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("The %s field max length is %s", fe.Field(), fe.Param())
	default:
		return fe.Error()
	}
}

func MakeValidationErrorMessage(err error) (string, bool) {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) && len(validationErrs) > 0 {
		return getValidationMessage(validationErrs[0]), true
	}

	return "", false
}
