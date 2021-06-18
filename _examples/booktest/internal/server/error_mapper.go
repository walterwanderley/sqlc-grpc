package server

import (
	"context"
	"database/sql"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"booktest/internal/validation"
)

func errorMapper(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	res, err := handler(ctx, req)
	if err != nil {
		if errors.Is(err, validation.ErrUserInput) {
			err = status.Error(codes.InvalidArgument, err.Error())
		} else if errors.Is(err, sql.ErrNoRows) {
			err = status.Error(codes.NotFound, err.Error())
		}
	}

	return res, err
}
