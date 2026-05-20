package models

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// posts model
type Post struct {
	ID      int `json:"id" db:"id"`
	Title   string `json:"title"  db:"title" validate:"min=3,max=30,required"`
	Content string `json:"content" db:"content" validate:"min=5,required"`
	UserID  int    `json:"user_id" db:"user_id" validate:"required"`
}

//likes model
type Like struct{
    Id  int `json:"id" db:"id"`
	UserID int `json:"user_id" db:"user_id" validate:"required,gt=0"`
	PostID int `json:"post_id"  db:"post_id" validate:"required,gt=0"`

	

}
var validate = validator.New()

func (post Post) Validate()string{
	var errMSGs []string 

	if err := validate.Struct(post); err != nil{
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

