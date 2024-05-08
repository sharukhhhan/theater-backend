package validation

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidatePayload(payload interface{}) error {
	validate := validator.New()
	err := validate.RegisterValidation("twoDecimalPlaces", twoDecimalPlaces)
	if err != nil {
		return err
	}

	var field, tag string
	err = validate.Struct(payload)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			field = e.Field()
			tag = e.Tag()
			return errors.New(fmt.Sprintf("validation error with %s: %s", field, tag))
		}
	}

	return nil
}

func twoDecimalPlaces(fl validator.FieldLevel) bool {
	// Use regex to match at most two decimal places
	re := regexp.MustCompile(`^\d+(\.\d{1,2})?$`)
	return re.MatchString(fmt.Sprintf("%v", fl.Field().Float()))
}
