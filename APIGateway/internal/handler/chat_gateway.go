package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/abelmalu/golang-posts/APIGateway/config"
	ierrors "github.com/abelmalu/golang-posts/APIGateway/internal/errors"
	"github.com/abelmalu/golang-posts/APIGateway/pkg/utils"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type ChatHandler struct {
	logger *platform.Logger
	config *config.Config
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var cfg *config.Config
var err error
var ChatProxy *httputil.ReverseProxy
var chatTarget *url.URL
var targetURL string

func init() {

	// load environment variables using godoenv package
	if err := godotenv.Load(); err != nil {

		log.Fatalf("Error while loading environment variables %v", err)

	}

	cfg, err = config.LoadConfig()

	if err != nil {

		log.Printf("Error loading env variables %v", err)
	}

	targetURL = "http://" + cfg.ChatServiceADD

	fmt.Println(targetURL)
	chatTarget, _ = url.Parse(targetURL)
	ChatProxy = httputil.NewSingleHostReverseProxy(chatTarget)

}

func NewChatHandler(l *platform.Logger) *ChatHandler {

	return &ChatHandler{

		logger: l,
		config: cfg,
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

	targetEndpoint := fmt.Sprintf("ws://%v/api/chat/ws", h.config.ChatServiceADD)

	fmt.Println(targetEndpoint,"target")

	backendConn, _, err := websocket.DefaultDialer.Dial(
		targetEndpoint,
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
	limit := c.Query("X-Limit")
	lastSeenID := c.Query("X-Last-ID")

	c.Request.Header.Set("X-Request-ID", requestID)
	c.Request.Header.Set("X-User-ID", userIDStr)
	c.Request.Header.Set("X-Limit", limit)
	c.Request.Header.Set("X-Last-ID", lastSeenID)

	h.logger.Info("proxying the request to notification service")
	ChatProxy.ServeHTTP(c.Writer, c.Request)

}

func (h *ChatHandler) GetChatMessages(c *gin.Context) {

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
	limit := c.Query("X-Limit")
	lastSeenID := c.Query("X-Last-ID")
	chatID := c.Param("id")

	c.Request.Header.Set("X-Request-ID", requestID)
	c.Request.Header.Set("X-User-ID", userIDStr)
	c.Request.Header.Set("X-Limit", limit)
	c.Request.Header.Set("X-Last-ID", lastSeenID)
	c.Request.Header.Set("Chat-ID", chatID)

	h.logger.Info("proxying the request to chat service")
	ChatProxy.ServeHTTP(c.Writer, c.Request)

}
