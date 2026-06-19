package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/abelmalu/golang-posts/Auth/internal/core"
	"github.com/abelmalu/golang-posts/Auth/internal/dto"
	ierrors "github.com/abelmalu/golang-posts/Auth/internal/errors"
	model "github.com/abelmalu/golang-posts/Auth/internal/models"
	"github.com/abelmalu/golang-posts/Auth/pkg/utils"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/metadata"
)

type AuthService struct {
	repo   core.AuthRepository
	logger *platform.Logger
	redis  *redis.Client
	kafka  sarama.SyncProducer
}

var maxLoginAttempt = 6
var tempMaxLoginAttempt = 3

func NewAuthService(authRepo core.AuthRepository, redisCient *redis.Client, kafkaClient sarama.SyncProducer, logger *platform.Logger) *AuthService {

	return &AuthService{
		repo:   authRepo,
		redis:  redisCient,
		kafka:  kafkaClient,
		logger: logger,
	}
}
func (authSer *AuthService) Register(ctx context.Context, user *model.User) (*model.User, *model.TokenPair, error) {
	var clientMetadata string
	var clientType model.ClientType
	if user.Name == "" {

		return nil, nil, ierrors.NewValidationError(ierrors.MSGNameIsRequired, nil, nil)
	}

	if user.Username == "" {

		return nil, nil, ierrors.NewValidationError(ierrors.MSGUsernameIsRequired, nil, nil)
	}

	if user.Password == "" {

		return nil, nil, ierrors.NewValidationError(ierrors.MSGPasswordIsRequired, nil, nil)
	}

	if user.Email == "" {

		return nil, nil, ierrors.NewValidationError(ierrors.MSGEmailIsRequired, nil, nil)
	}
	createdUser, err := authSer.repo.Register(ctx, user)

	if err != nil {

		return nil, nil, err
	}

	createdUserByte, err := json.Marshal(createdUser)

	if err != nil {

		return nil, nil, ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, err)
	}
	msg := &sarama.ProducerMessage{
		Topic: "userCreated",
		Value: sarama.StringEncoder(createdUserByte),
	}

	// Send the message to Kafka
	partition, offset, err := authSer.kafka.SendMessage(msg)

	authSer.logger.Info(fmt.Sprintf("Apache kafka Partition : %d Offset : %d", partition, offset))

	md, exists := metadata.FromIncomingContext(ctx)

	if !exists {
		return nil, nil, ierrors.NewValidationError(ierrors.MSGUnkownDevice, nil, nil)
	}
	values := md.Get("x-client-type")
	if len(values) > 0 {
		clientMetadata = values[0]
	} else {

		return nil, nil, ierrors.NewValidationError(ierrors.MSGUnkownDevice, nil, nil)

	}
	switch clientMetadata {
	case "web":
		clientType = model.ClientWeb

	case "mobile":
		clientType = model.ClientMobile
	default:

		return nil, nil, ierrors.NewValidationError(ierrors.MSGUnkownDevice, nil, nil)

	}
	tokens, err := authSer.issueTokens(createdUser.ID, clientType, createdUser.Role)
	if err != nil {

		return nil, nil, err
	}

	return createdUser, tokens, nil
}

func (authSer *AuthService) Login(ctx context.Context, userName, password string) (*model.User, *model.TokenPair, error) {
	var clientMetadata string
	var clientType model.ClientType
	if userName == "" {

		return nil, nil, ierrors.NewValidationError(ierrors.MSGUsernameIsRequired, nil, nil)

	}

	//checking if the user canceled the request or the reques timeouts
	if err := ctx.Err(); err != nil {

		switch {
		case errors.Is(err, context.Canceled):

			return nil, nil, ierrors.NewCancelationError(ierrors.MSGRequestCanceled, err)

		case errors.Is(err, context.DeadlineExceeded):

			return nil, nil, ierrors.NewTimeoutError(ierrors.MSGTimeoutReached, err)

		}

	}

	if password == "" {

		return nil, nil, ierrors.NewValidationError(ierrors.MSGPasswordIsRequired, nil, nil)
	}

	fetchedUser, err := authSer.repo.Login(ctx, userName, password)
	if err != nil {
		return nil, nil, err

	}

	// if user is found check the password
	if fetchedUser.Password != password {

		fetchedUser.FailedLoginAttempts++

		if fetchedUser.FailedLoginAttempts == tempMaxLoginAttempt {

			lockUntil := (time.Now().Add(time.Minute * 3))

			fetchedUser.TemporaryLockUntil = &lockUntil

			_, err := authSer.repo.TemporaryLockUntil(ctx, fetchedUser)
			if err != nil {

				return nil, nil, err
			}

			return nil, nil, ierrors.NewUnauthorizedError(ierrors.ErrorMessage("Your Account is temporarly blocked"), nil)

		}

		if fetchedUser.FailedLoginAttempts >= maxLoginAttempt {
			_, err := authSer.repo.UpdateFailedLoginAttempst(ctx, fetchedUser)
			if err != nil {

				return nil, nil, err
			}

			return nil, nil, ierrors.NewUnauthorizedError(ierrors.ErrorMessage("Your Account is Prmanently Blocked Contact Adminstrator"), nil)

		}

		_, err := authSer.repo.UpdateFailedLoginAttempst(ctx, fetchedUser)
		if err != nil {

			return nil, nil, err
		}
		
		remainingAttempts := maxLoginAttempt - fetchedUser.FailedLoginAttempts

		return nil, nil, ierrors.NewUnauthorizedError(ierrors.ErrorMessage(fmt.Sprintf("Invalid Credentials %d Attempts left",remainingAttempts)), nil)

	}

	md, exists := metadata.FromIncomingContext(ctx)

	if !exists {
		return nil, nil, ierrors.NewValidationError(ierrors.MSGUnkownDevice, nil, nil)
	}
	values := md.Get("x-client-type")
	if len(values) > 0 {
		clientMetadata = values[0]
	} else {

		return nil, nil, ierrors.NewValidationError(ierrors.MSGUnkownDevice, nil, nil)

	}
	switch clientMetadata {
	case "web":
		clientType = model.ClientWeb

	case "mobile":
		clientType = model.ClientMobile
	default:

		return nil, nil, ierrors.NewValidationError(ierrors.MSGUnkownDevice, nil, nil)

	}
	tokens, err := authSer.issueTokens(fetchedUser.ID, clientType, fetchedUser.Role)
	if err != nil {

		return nil, nil, err
	}

	return fetchedUser, tokens, nil

}

func (authSer *AuthService) Logout(ctx context.Context, refreshToken string) error {

	// validate the token to check if it is tampered
	_, err := utils.ValidateRefreshToken(refreshToken)

	if err != nil {

		return ierrors.NewUnauthorizedError(ierrors.MSGFailedToValidateToken, nil)

	}
	JTI, err := utils.GetJTI(ctx)
	if err != nil {

		return err
	}
	expTime, err := utils.GetJWTEXPTime(ctx)
	if err != nil {

		return err
	}
	expUnix, err := strconv.Atoi(expTime)
	if err != nil {

		return err
	}

	expirationTime := time.Unix(int64(expUnix), 0)

	expDuration := time.Until(expirationTime)
	authSer.redis.Set(ctx, JTI, 1, expDuration)

	// hash the token to check with DB token
	hashedRefreshToken := utils.HashToken(refreshToken)

	if err := authSer.repo.Logout(ctx, hashedRefreshToken); err != nil {

		return err

	}

	return err
}

func (authSer *AuthService) RefreshHandler(ctx context.Context, refreshToken string) (*model.TokenPair, error) {

	if refreshToken == "" {

		return nil, ierrors.NewValidationError(ierrors.MSGRefreshTokenIsRequired, nil, nil)
	}

	// Get the refresh token from the DB
	tokenRecord, err := authSer.repo.GetRefreshToken(refreshToken)
	if err != nil {

		return nil, err
	}
	// check if it is revoked or has expired
	if tokenRecord.Revoked || tokenRecord.ExpiresAt.Before(time.Now()) {
		return nil, ierrors.NewUnauthorizedError(ierrors.MSGUnauthorizedAccess, nil)

	}
	// validate the token to check if it is tampered or expired
	_, err = utils.ValidateRefreshToken(refreshToken)

	if err != nil {

		return nil, ierrors.NewUnauthorizedError(ierrors.MSGUnauthorizedAccess, nil)

	}

	user, err := authSer.repo.GetUserByID(tokenRecord.ID)

	if err != nil {
		return nil, err

	}
	var clientType model.ClientType
	userID := tokenRecord.UserID
	clientTypeStr := tokenRecord.ClientType

	switch clientTypeStr {
	case "web":
		clientType = model.ClientWeb
	case "mobile":
		clientType = model.ClientMobile

	}

	// revoke the old token so it can't be used anymore
	if err := authSer.repo.RevokeRefreshToken(tokenRecord.TokenText); err != nil {

		return nil, err
	}

	// Generate new tokens Rotate refresh token (issue a new one)
	tokens, err := authSer.issueTokens(userID, clientType, user.Role)
	if err != nil {

		return nil, err
	}

	// store the refresh token
	newExpireTime := time.Now().Add(24 * 30 * time.Hour)
	_, err = authSer.repo.StoreRefreshTokens(userID, tokens.RefreshToken, newExpireTime, clientTypeStr)
	if err != nil {

		return nil, err

	}

	return tokens, nil

}

func (authSer *AuthService) issueTokens(userID int, clientType model.ClientType, userRole string) (*model.TokenPair, error) {

	accessToken, err := utils.GenerateAcessToken(userID, userRole)
	if err != nil {

		return nil, ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, err)
	}
	refreshToken, err, expiresAt := utils.GenerateRefreshToken(userID)
	if err != nil {

		return nil, ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, err)

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

func (authSer *AuthService) SearchUser(ctx context.Context, username, cursor string, limit int) (*dto.PaginatedResponse, error) {

	searchPattern := "%" + username + "%"

	resp, err := authSer.repo.SearchUser(ctx, searchPattern, cursor, limit)

	if err != nil {

		return nil, err
	}

	return resp, nil
}
