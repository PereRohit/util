package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func Validate(s interface{}) error {
	err := validator.New().Struct(s)
	if errs, ok := err.(validator.ValidationErrors); ok {
		validationErrs := "validation "
		for i, err := range errs {
			validationErrs += fmt.Sprintf("%d: field <%s> with value <%s> failed for <%s> validation.\n",
				i+1, err.Field(), err.Value().(string), err.Tag())
		}
		return fmt.Errorf(validationErrs)
	}
	return err
}
