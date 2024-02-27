package validate

import (
	"errors"
	"github.com/jason810496/Dcard-Advertisement-API/pkg/schemas"
)

// go-playground/validator/v10 can't validate omitempty,gte=1 when the value is 0
func PublicAdRequestAge(json *schemas.PublicAdRequest) error {
	if json.Age <= 0 {
		return errors.New("Age should be greater than 0")
	}

	return nil
}
