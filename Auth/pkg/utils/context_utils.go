package utils

import (
	"context"

	ierrors "github.com/abelmalu/golang-posts/Auth/internal/errors"
	"google.golang.org/grpc/metadata"
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

func GetJTI(c context.Context) (string, error) {

	var JTI string

	md, exists := metadata.FromIncomingContext(c)

	if !exists {
		return "", ierrors.ErrMetaDataNotFound
	}
	values := md.Get("JTI")
	if len(values) > 0 {
		JTI = values[0]
	} else {

		return "", ierrors.ErrJTINotFound

	}
	return JTI, nil

}


func GetJWTEXPTime(c context.Context) (string, error) {

	var expTime string

	md, exists := metadata.FromIncomingContext(c)

	if !exists {
		return "", ierrors.ErrMetaDataNotFound
	}
	values := md.Get("exp")
	if len(values) > 0 {
		expTime = values[0]
	} else {

		return "", ierrors.ErrJTINotFound

	}
	return expTime, nil

}
