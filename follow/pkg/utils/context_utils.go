package utils

import (
	"context"
	"strconv"

	ierrors "github.com/abelmalu/golang-posts/follow/internal/errors"
	"google.golang.org/grpc/metadata"
)

const (
	MaxLimit     = 100
	DefaultLimit = 20
)

func GetRequestID(c context.Context) (string, error) {
	var requestID string

	md, exists := metadata.FromIncomingContext(c)

	if !exists {
		return "", ierrors.ErrMetaDataNotFound
	}
	values := md.Get("request-id")
	if len(values) > 0 {
		requestID = values[0]
	} else {

		return "", ierrors.ErrRequestIDNotFound

	}
	return requestID, nil

}

func GetUserID(c context.Context) (int, error) {
	var userIDStr string

	md, exists := metadata.FromIncomingContext(c)

	if !exists {
		return 0, ierrors.ErrMetaDataNotFound
	}
	values := md.Get("user-id")
	if len(values) > 0 {
		userIDStr = values[0]
	} else {

		return 0, ierrors.ErrUSerIDNotFound

	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {

		return 0, ierrors.ErrUSerIDNotFound

	}
	return userID, nil

}

func GetPostID(c context.Context) (int, error) {
	var postIDStr string

	md, exists := metadata.FromIncomingContext(c)

	if !exists {
		return 0, ierrors.ErrMetaDataNotFound
	}
	values := md.Get("post-id")
	if len(values) > 0 {
		postIDStr = values[0]
	} else {

		return 0, ierrors.ErrPostIDNotFound

	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {

		return 0, ierrors.ErrPostIDNotFound

	}
	return postID, nil

}

func ValidatePaginationLimit(limit int) int {

	if limit <= 0 {

		return DefaultLimit

	}

	if limit > 100 {

		return MaxLimit
	}

	return limit
}
