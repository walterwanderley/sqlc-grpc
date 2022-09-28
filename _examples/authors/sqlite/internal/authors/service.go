// Code generated by sqlc-grpc (https://github.com/walterwanderley/sqlc-grpc). DO NOT EDIT.

package authors

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	pb "authors/api/authors/v1"
)

type Service struct {
	pb.UnimplementedAuthorsServiceServer
	logger  *zap.Logger
	querier *Queries
}

func (s *Service) CreateAuthor(ctx context.Context, req *pb.CreateAuthorRequest) (*pb.CreateAuthorResponse, error) {
	var arg CreateAuthorParams
	arg.Name = req.GetName()
	if v := req.GetBio(); v != nil {
		arg.Bio = sql.NullString{Valid: true, String: v.Value}
	}

	result, err := s.querier.CreateAuthor(ctx, arg)
	if err != nil {
		s.logger.Error("CreateAuthor sql call failed", zap.Error(err))
		return nil, err
	}
	return &pb.CreateAuthorResponse{Value: toExecResult(result)}, nil
}

func (s *Service) DeleteAuthor(ctx context.Context, req *pb.DeleteAuthorRequest) (*pb.DeleteAuthorResponse, error) {
	id := req.GetId()

	err := s.querier.DeleteAuthor(ctx, id)
	if err != nil {
		s.logger.Error("DeleteAuthor sql call failed", zap.Error(err))
		return nil, err
	}
	return &pb.DeleteAuthorResponse{}, nil
}

func (s *Service) GetAuthor(ctx context.Context, req *pb.GetAuthorRequest) (*pb.GetAuthorResponse, error) {
	id := req.GetId()

	result, err := s.querier.GetAuthor(ctx, id)
	if err != nil {
		s.logger.Error("GetAuthor sql call failed", zap.Error(err))
		return nil, err
	}
	return &pb.GetAuthorResponse{Author: toAuthor(result)}, nil
}

func (s *Service) ListAuthors(ctx context.Context, req *pb.ListAuthorsRequest) (*pb.ListAuthorsResponse, error) {

	result, err := s.querier.ListAuthors(ctx)
	if err != nil {
		s.logger.Error("ListAuthors sql call failed", zap.Error(err))
		return nil, err
	}
	res := new(pb.ListAuthorsResponse)
	for _, r := range result {
		res.List = append(res.List, toAuthor(r))
	}
	return res, nil
}
