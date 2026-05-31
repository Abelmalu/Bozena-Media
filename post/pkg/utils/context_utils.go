package utils

import (
	"context"

	ierrors "github.com/abelmalu/golang-posts/post/internal/errors"
	"google.golang.org/grpc/metadata"
)
const (
	MaxLimit     = 100
	DefaultLimit = 20
)


func GetRequestID(c context.Context) (string,error){
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
	return requestID,nil


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