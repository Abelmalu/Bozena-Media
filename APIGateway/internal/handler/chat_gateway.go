package handler

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	ierrors "github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	"github.com/abelmalu/golang-posts/APIGateway/pkg/utils"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type ChatHandler struct {
	logger *platform.Logger
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var chatTarget, _ = url.Parse("http://localhost:8084")

var ChatProxy = httputil.NewSingleHostReverseProxy(target)

func NewChatHandler(l *platform.Logger) *ChatHandler {

	return &ChatHandler{

		logger: l,
	}

}

func (h *ChatHandler) Connect(c *gin.Context) {
	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			h.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			h.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

	}

	userID, err := utils.GetUserID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrUserIDNotFoundInContext) {

			h.logger.Error("couldn't couldn't find userID in the context", zap.String("type", "something went wrong"), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			h.logger.Error("couldn't assert the user ID to string", zap.String("type", "something went wrong"), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}

	}

	userIDStr := strconv.Itoa(userID)

	c.Request.Header.Set("X-Request-ID", requestID)
	c.Request.Header.Set("X-User-ID", userIDStr)

	clientConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer clientConn.Close()

	backendConn, _, err := websocket.DefaultDialer.Dial(
		"ws://localhost:8084/api/chat/ws",
		http.Header{
			"X-User-ID":    []string{userIDStr},
			"X-Request-ID": []string{requestID},
		},
	)
	if err != nil {
		return
	}
	defer backendConn.Close()

	go func() {
		for {
			mt, msg, err := clientConn.ReadMessage()
			if err != nil {
				backendConn.Close()
				return
			}

			if err := backendConn.WriteMessage(mt, msg); err != nil {
				return
			}
		}
	}()

	for {
		mt, msg, err := backendConn.ReadMessage()
		if err != nil {
			break
		}

		if err := clientConn.WriteMessage(mt, msg); err != nil {
			break
		}
	}
}

func (h *ChatHandler) GetUserChats(c *gin.Context) {

	requestID, err := utils.GetRequestID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrRequestIDNotFoundInContext) {

			h.logger.Error("couldn't get request ID from context", zap.Error(errors.New("couldn't find request ID")))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))

			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			h.logger.Error("couldn't assert the request ID to string", zap.String("type", "something went wrong"))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return
		}

	}

	userID, err := utils.GetUserID(c)
	if err != nil {

		if errors.Is(err, ierrors.ErrUserIDNotFoundInContext) {

			h.logger.Error("couldn't couldn't find userID in the context", zap.String("type", "something went wrong"), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}
		if errors.Is(err, ierrors.ErrTypeAssertionFailed) {

			h.logger.Error("couldn't assert the user ID to string", zap.String("type", "something went wrong"), zap.String("requestID", requestID))
			c.Error(ierrors.NewInternalError(ierrors.MSGSomethingWentWrong, nil))
			return

		}

	}

	userIDStr := strconv.Itoa(userID)

	c.Request.Header.Set("X-Request-ID", requestID)
	c.Request.Header.Set("X-User-ID", userIDStr)

	h.logger.Info("proxying the request to notification service")
	ChatProxy.ServeHTTP(c.Writer, c.Request)

}
