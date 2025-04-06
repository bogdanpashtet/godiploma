package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorFromDomain(err error) error {
	switch status.Code(err) {
	case codes.InvalidArgument:
		return status.Error(codes.InvalidArgument, err.Error())
	case codes.Canceled:
		return status.Error(codes.Canceled, err.Error())
	case codes.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, err.Error())
	}

	return status.Error(codes.Internal, err.Error())
}
