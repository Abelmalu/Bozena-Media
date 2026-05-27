package middleware

import (
	"errors"
	"strings"

	ierrors "github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	"github.com/abelmalu/golang-posts/APIGateway/pkg/utils"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func AuthMiddleware(logger *platform.Logger, redisclient *redis.Client) gin.HandlerFunc {
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
				return
			}

		}
		authHeader := c.GetHeader("Authorization")

		if strings.HasPrefix(authHeader, "Bearer ") {

			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if tokenStr == "" {

			err := ierrors.NewValidationError(ierrors.MSGUnauthorizedAccess, nil, nil)

			logger.Error("token string not found in authorization header",zap.String("requestID",requestID))

			c.Error(err)
			c.Abort()
			return
		}

		// Validate the token

		tokenClaims, err := utils.ValidateAccessToken(tokenStr)

		if err != nil {
			err := ierrors.NewValidationError(ierrors.MSGUnauthorizedAccess, nil, nil)

			logger.Error("token string not found in authorization header")

			c.Error(err)
			c.Abort()
			return
		}

		uid, ok := tokenClaims["user_id"].(float64)
		if !ok {
			err := ierrors.NewValidationError(ierrors.MSGUnauthorizedAccess, nil, nil)

			logger.Error("Invalid token claims:Failed while asserting user_id", zap.Error(err))

			c.Error(err)
			c.Abort()
			return
		}
		c.Set("userID", int(uid))
		userRole, ok := tokenClaims["userRole"].(string)
		if !ok {
			err := ierrors.NewValidationError(ierrors.MSGUnauthorizedAccess, nil, nil)

			logger.Error("Invalid token claims:Failed while asserting userRole", zap.Error(err))

			c.Error(err)
			c.Abort()
			return
		}
		c.Set("userRole", userRole)

		JTI, ok := tokenClaims["jti"].(string)
		if !ok {
			err := ierrors.NewValidationError(ierrors.MSGUnauthorizedAccess, nil, nil)

			logger.Error("Invalid token claims:Failed while asserting JTI", zap.Error(err))

			c.Error(err)
			c.Abort()
			return
		}
		blackListedtoken, err := redisclient.Get(c.Request.Context(), JTI).Result()

		if err == nil {
			internalErr := ierrors.NewValidationError(ierrors.MSGUnauthorizedAccess, nil, nil)
			logger.Warn("logged out token reuse",zap.String("JTI",blackListedtoken))
			c.Error(internalErr)
			c.Abort()
			return

		}

		c.Set("JTI", JTI)

		expTime, ok := tokenClaims["exp"].(float64)
		if !ok {
			err := ierrors.NewValidationError(ierrors.MSGUnauthorizedAccess, nil, nil)

			logger.Error("Invalid token claims:Failed while asserting expTime", zap.Error(err))

			c.Error(err)
			c.Abort()
			return
		}

		c.Set("expTime", expTime)

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
