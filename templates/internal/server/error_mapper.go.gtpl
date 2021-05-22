package server

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrUserInput = errors.New("")

func errorMapper(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	res, err := handler(ctx, req)
	if err != nil {
		if errors.Is(err, ErrUserInput) {
			err = status.Error(codes.InvalidArgument, err.Error())
		} else if err.Error() == "sql: no rows in result set" {
			err = status.Error(codes.NotFound, err.Error())
		}
	}

	return res, err
}
