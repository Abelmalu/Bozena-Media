package dto

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ToggleLikeRequest struct {
	State bool `json:"state" validate:"required"`
}

type ToggleLikeResponse struct {
	Message string `json:"message"`
}

var validate = validator.New()

func (toggleLikeRequest ToggleLikeRequest) Validate() string {
	var errMSGs []string

	if err := validate.Struct(toggleLikeRequest); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {

			for _, fieldErr := range validationErrors {

				switch fieldErr.Tag() {

				case "required":
					errMSGs = append(errMSGs, fmt.Sprintf("%s is required", fieldErr.Field()))
				
				

				}

			}

			errMSGsStr := strings.Join(errMSGs, ",")

			return errMSGsStr

		}

	}
	return ""

}
