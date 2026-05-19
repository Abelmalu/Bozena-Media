package model

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// users model
type User struct{

	ID int `json:"id" db:"id" `
	Name string `json:"name" db:"name" validate:"required,min=2,max=30"`
	Username string `json:"username" db:"username" validate:"required,min=2,max=30"`
	Email string `json:"email" db:"email" validate:"email,required"`
	Password string `json:"password" db:"password" validate:"required,min=8,max=30"`
	Role string `json:"role" db:"role"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}


// client types for mobile apps and browsers
type ClientType string 

const (
	ClientWeb ClientType = "web"
	ClientMobile ClientType = "mobile"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
var validate = validator.New()

func (user User) Validate()string{
	var errMSGs []string 

	if err := validate.Struct(user); err != nil{
		if validationErrors,ok := err.(validator.ValidationErrors);ok{

			

			for _,fieldErr := range validationErrors{

				switch fieldErr.Tag(){

				case "required":
				errMSGs = append(errMSGs, fmt.Sprintf("%s is required", fieldErr.Field()))
				case "min":
				errMSGs = append(errMSGs, fmt.Sprintf("%s must be atleast  %s characters", fieldErr.Field(), fieldErr.Param()))
					
				case "max":
					errMSGs = append(errMSGs, fmt.Sprintf("%s must be below %s characters",fieldErr.Field(),fieldErr.Param()))
				case "email":
					errMSGs = append(errMSGs, "Please enter a valid email address")



				}

			}

			errMSGsStr := strings.Join(errMSGs,",")

			return errMSGsStr


		}


	}
	return ""

	



}
