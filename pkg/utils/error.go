package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Error handling for `validate.PublicAdRequestAge`
// ( go-playground/validator/v10 can't validate omitempty,gte=1 when the value is 0 )
func NewAgeError(ctx *gin.Context, status int) {
	ctx.AbortWithStatusJSON(status, HTTPError{
		Code: status,
		Errors: []ErrorMsg{
			{
				Field:   "Age",
				Message: "Should be greater than 1",
			},
		},
	})
}

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
	case "gte", "gtefield", "min":
		return "Should be greater than " + fe.Param()
	case "oneof":
		return "Should be one of " + fe.Param()
	case "max":
		return "Should be less than " + fe.Param()
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
