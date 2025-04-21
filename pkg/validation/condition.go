package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	tempRe = regexp.MustCompile(`^temp\s*(<=|>=|<|>|==|=|!=)\s*[0-9]+(?:\.[0-9]+)?$`)
	condRe = regexp.MustCompile(`^condition\s*=\s*[A-Za-z]+$`)
)

func RegisterConditionValidator(v *validator.Validate) {
	v.RegisterValidation("condition", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		return tempRe.MatchString(s) || condRe.MatchString(s)
	})
}
