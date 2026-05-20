package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// V wraps go-playground/validator with a simple API.
type V struct {
	v *validator.Validate
}

func New() *V {
	return &V{v: validator.New()}
}

func (v *V) ValidateStruct(data interface{}) map[string]string {
	err := v.v.Struct(data)
	if err == nil {
		return nil
	}
	out := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		out[strings.ToLower(e.Field())] = fieldError(e)
	}
	return out
}

func fieldError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "required"
	case "email":
		return "invalid email"
	case "min":
		return fmt.Sprintf("must be at least %s characters", e.Param())
	default:
		return e.Tag()
	}
}
