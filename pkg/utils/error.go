package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// NewError example
func NewError(ctx *gin.Context, status int, err error) {
	var validate_errors validator.ValidationErrors

	if errors.As(err, &validate_errors) {
		out := make([]ErrorMsg, len(validate_errors))
		for i, err := range validate_errors {
			out[i] = ErrorMsg{
				Field:   err.Field(),
				Message: getErrorMsg(err),
			}
		}

		ctx.AbortWithStatusJSON(status, HTTPError{
			Code:   status,
			Errors: out,
		})
	}
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "gte", "gtefield":
		return "Should be greater than " + fe.Param()
	case "oneof":
		return "Should be one of " + fe.Param()
	}
	return "Unknown error"
}

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message" example:"bad request"`
}

// HTTPError example
type HTTPError struct {
	Code   int        `json:"code" example:"400"`
	Errors []ErrorMsg `json:"errors"`
}
