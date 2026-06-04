package dto

import (
	"fmt"
	"strings"

	"github.com/abelmalu/golang-posts/follow/internal/models"
	"github.com/go-playground/validator/v10"
)



type FollowRequest struct {

	Follow *bool `validate:"required"`


}

type FollowResponse struct {

	Message string 


}


type PaginatedFollowersResponse struct {

	Followers []*models.UserFollowers
	Cursor string 
	HasNext bool
}


type PaginatedFollowingsResponse struct {

	Followings []*models.Follow
	Cursor string 
	HasNext bool
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