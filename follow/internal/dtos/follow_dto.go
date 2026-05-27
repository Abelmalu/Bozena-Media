package dto

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)



type FollowRequest struct {

	Follow *bool `validate:"required"`


}

type FollowResponse struct {

	Message string 


}

var validate = validator.New()


func (followRequest *FollowRequest) Validate() string {

	var errMSGs []string

	if err := validate.Struct(followRequest); err != nil {
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