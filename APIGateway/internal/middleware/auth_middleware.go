package middleware

import (
	"errors"
	"strings"

	"github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	"github.com/abelmalu/golang-posts/APIGateway/pkg/utils"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthMiddleware(logger *platform.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string
			requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

		logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

		}

	

		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

		}
		if errors.Is(err, ierrors.ErrUserIDNotFoundInContext) {

		}
	}
		authHeader := c.GetHeader("Authorization")

		if strings.HasPrefix(authHeader, "Bearer ") {

			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if tokenStr == "" {

			err := ierrors.NewValidationError(ierrors.MSGUnauthorizedAccess, nil, nil)

			logger.Error("token string not found in authorization header")

			utils.SendErrorResponse[error](c, err, requestID, err.HTTPStatus())

			c.Abort()
			return

		}

		// Validate the token

		tokenClaims, err := utils.ValidateAccessToken(tokenStr)

		if err != nil {
			err := ierrors.NewValidationError(ierrors.MSGUnauthorizedAccess, nil, nil)

			logger.Error("token string not found in authorization header")

			utils.SendErrorResponse[error](c, err, requestID, err.HTTPStatus())
			c.Abort()
			return
		}

		uid, ok := tokenClaims["user_id"].(float64)
		if !ok {
			err := ierrors.NewValidationError(ierrors.MSGUnauthorizedAccess, nil, nil)

			logger.Error("Invalid token claims:Failed while asserting user_id", zap.Error(err))

			utils.SendErrorResponse[error](c, err, requestID, err.HTTPStatus())

			return
		}
		c.Set("userID", int(uid))
		userRole, ok := tokenClaims["userRole"].(string)
		if !ok {
			err := ierrors.NewValidationError(ierrors.MSGUnauthorizedAccess, nil, nil)

			logger.Error("Invalid token claims:Failed while asserting userRole", zap.Error(err))

			utils.SendErrorResponse[error](c, err, requestID, err.HTTPStatus())

			return
		}
		c.Set("userRole", userRole)
		// convert once here

		c.Next() // Token is valid, proceed to the next handler!
	}
}

func AuthorizeRoles(logger *platform.Logger, allowedRoles ...string) gin.HandlerFunc {

	return func(c *gin.Context) {
			requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			c.Abort()
			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			c.Abort()
			return

		}

	

		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

		}
		if errors.Is(err, ierrors.ErrUserIDNotFoundInContext) {

		}
	}


		role, ok := c.Get("userRole")
		if !ok {
			err := ierrors.NewValidationError(ierrors.MSGUnauthorizedAccess, nil, nil)

			logger.Error("Invalid token claims:Failed while asserting userRole", zap.Error(err))

			utils.SendErrorResponse[error](c, err, requestID, err.HTTPStatus())

			c.Abort()
			return
		}
		userRole := role.(string)

		hasAccess := false
		for _, r := range allowedRoles {
			if r == userRole {
				hasAccess = true
				break

			}

		}
		if !hasAccess {

			err := ierrors.NewValidationError(ierrors.MSGUnauthorizedAccess, nil, nil)

			logger.Warn("Invalid token claims:Failed while asserting userRole", zap.Error(err))

			utils.SendErrorResponse[error](c, err, requestID, err.HTTPStatus())

			c.Abort()

		}
		c.Next()

	}
}
