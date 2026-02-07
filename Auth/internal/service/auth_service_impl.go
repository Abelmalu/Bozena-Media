package service

import (
	"context"
	"errors"

	"github.com/abelmalu/golang-posts/Auth/internal/core"
	model "github.com/abelmalu/golang-posts/Auth/internal/models"
	"github.com/abelmalu/golang-posts/pkg"
	"google.golang.org/grpc/metadata"
)

type AuthService struct {
	repo core.AuthRepository
}


func NewAuthService(authRepo core.AuthRepository) *AuthService {

	return &AuthService{repo: authRepo}
}
func (authSer *AuthService) Register(ctx context.Context, user *model.User) (*model.User,*model.TokenPair, error) {
	var clientMetadata string
	var clientType model.ClientType
	if user.Name == "" {

		return nil,nil, errors.New("name is required")
	}

	if user.Username == "" {

		return nil, nil, errors.New("username is required")
	}

	if user.Password == "" {

		return nil,nil, errors.New("password is required")
	}

	if user.Email == "" {

		return nil, nil, errors.New("email is required")
	}
	createdUser, err := authSer.repo.Register(ctx, user)

	if err != nil {

		return nil, nil,err
	}
	md, exists := metadata.FromIncomingContext(ctx)

	if !exists {
		return nil,nil, errors.New("Unknown device type")
	}
	values := md.Get("x-client-type")
	if len(values) > 0 {
		clientMetadata = values[0]
	} else {

		return nil, nil,errors.New("Unknown device type")

	}
	if clientMetadata == "web" {
		clientType = model.ClientWeb

	} else if clientMetadata == "mobile" {
		clientType = model.ClientMobile
	} else {

		return nil, nil,errors.New("Unknown device type")

	}
	tokens, err := authSer.issueTokens(createdUser.ID, clientType, createdUser.Role)

	return createdUser, tokens,nil
}

func (authSer *AuthService) Login(ctx context.Context,userName,password string)(*model.User,*model.TokenPair,error){
	var clientMetadata string
	var clientType model.ClientType
	if userName == "" {

		return nil, nil, errors.New("username is required")
	}

	if password == "" {

		return nil,nil, errors.New("password is required")
	}

	fetchedUser, err := authSer.repo.Login(ctx,userName,password)
	if err != nil {
		return nil,nil,errors.New("username already exists")
	}

	md, exists := metadata.FromIncomingContext(ctx)

	if !exists {
		return nil,nil, errors.New("Unknown device type")
	}
	values := md.Get("x-client-type")
	if len(values) > 0 {
		clientMetadata = values[0]
	} else {

		return nil, nil,errors.New("Unknown device type")

	}
	if clientMetadata == "web" {
		clientType = model.ClientWeb

	} else if clientMetadata == "mobile" {
		clientType = model.ClientMobile
	} else {

		return nil, nil,errors.New("Unknown device type")

	}
	tokens, err := authSer.issueTokens(fetchedUser.ID, clientType, fetchedUser.Role)

	return fetchedUser,tokens,nil






}

func (authSer *AuthService) issueTokens(userID int, clientType model.ClientType, userRole string) (*model.TokenPair, error) {

	accessToken, err := pkg.GenerateAcessToken(userID, userRole)
	if err != nil {

		return nil, err
	}
	refreshToken, err, expiresAt := pkg.GenerateRefreshToken(userID)
	if err != nil {

		return nil, err
	}

	_, err = authSer.repo.StoreRefreshTokens(userID, refreshToken, expiresAt, string(clientType))

	if err != nil {

		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

