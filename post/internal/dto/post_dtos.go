package dto

import (
	"fmt"
	"strings"

	"github.com/abelmalu/golang-posts/post/internal/models"
	"github.com/go-playground/validator/v10"
)


type UpdatePostRequest struct{
	Title   string `json:"title"  db:"title" validate:"min=3,max=30,required"`
	Content string `json:"content" db:"content" validate:"min=5,required"`
}

type PaginatedResponse struct {
	Posts   *[]models.Post
	Cursor  string
	HasNext bool
}


var validate = validator.New()

func (updatePost UpdatePostRequest) Validate()string{
	var errMSGs []string 

	if err := validate.Struct(updatePost); err != nil{
		if validationErrors,ok := err.(validator.ValidationErrors);ok{

			

			for _,fieldErr := range validationErrors {

				switch fieldErr.Tag(){

				case "required":
				errMSGs = append(errMSGs, fmt.Sprintf("%s is required", fieldErr.Field()))
				case "min":
				errMSGs = append(errMSGs, fmt.Sprintf("%s must be atleast  %s characters", fieldErr.Field(), fieldErr.Param()))
					
				case "max":
					errMSGs = append(errMSGs, fmt.Sprintf("%s must be below %s characters",fieldErr.Field(),fieldErr.Param()))
				
			

				}

			}

			errMSGsStr := strings.Join(errMSGs,",")

			return errMSGsStr


		}


	}
	return ""

	



}

