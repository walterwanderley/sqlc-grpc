package books

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "booktest/proto/books"
)

type Service struct {
	pb.UnimplementedBooksServiceServer
	logger *zap.Logger
	db     *Queries
}

func NewService(logger *zap.Logger, db *Queries) *Service {
	return &Service{logger: logger, db: db}
}

func (s *Service) BooksByTags(ctx context.Context, in *pb.BooksByTagsParams) (out *pb.BooksByTagsResponse, err error) {
	dollar_1 := in.GetDollar_1()

	result, err := s.db.BooksByTags(ctx, dollar_1)
	if err != nil {
		s.logger.Error("BooksByTags sql call failed", zap.Error(err))
		return
	}
	out = new(pb.BooksByTagsResponse)
	for _, r := range result {
		var item *pb.BooksByTagsRow
		item, err = toBooksByTagsRowProto(r)
		if err != nil {
			return
		}
		out.Value = append(out.Value, item)
	}
	return

}

func (s *Service) BooksByTitleYear(ctx context.Context, in *pb.BooksByTitleYearParams) (out *pb.BooksByTitleYearResponse, err error) {
	arg, err := toBooksByTitleYearParams(in)
	if err != nil {
		s.logger.Error("BooksByTitleYear input adapter failed", zap.Error(err))
		return
	}

	result, err := s.db.BooksByTitleYear(ctx, arg)
	if err != nil {
		s.logger.Error("BooksByTitleYear sql call failed", zap.Error(err))
		return
	}
	out = new(pb.BooksByTitleYearResponse)
	for _, r := range result {
		var item *pb.Book
		item, err = toBookProto(r)
		if err != nil {
			return
		}
		out.Value = append(out.Value, item)
	}
	return

}

func (s *Service) CreateAuthor(ctx context.Context, in *pb.CreateAuthorParams) (out *pb.Author, err error) {
	name := in.GetName()

	result, err := s.db.CreateAuthor(ctx, name)
	if err != nil {
		s.logger.Error("CreateAuthor sql call failed", zap.Error(err))
		return
	}
	return toAuthorProto(result)

}

func (s *Service) CreateBook(ctx context.Context, in *pb.CreateBookParams) (out *pb.Book, err error) {
	arg, err := toCreateBookParams(in)
	if err != nil {
		s.logger.Error("CreateBook input adapter failed", zap.Error(err))
		return
	}

	result, err := s.db.CreateBook(ctx, arg)
	if err != nil {
		s.logger.Error("CreateBook sql call failed", zap.Error(err))
		return
	}
	return toBookProto(result)

}

func (s *Service) DeleteBook(ctx context.Context, in *pb.DeleteBookParams) (out *emptypb.Empty, err error) {
	bookID := in.GetBookID()

	err = s.db.DeleteBook(ctx, bookID)
	if err != nil {
		s.logger.Error("DeleteBook sql call failed", zap.Error(err))
		return
	}
	return &emptypb.Empty{}, nil

}

func (s *Service) GetAuthor(ctx context.Context, in *pb.GetAuthorParams) (out *pb.Author, err error) {
	authorID := in.GetAuthorID()

	result, err := s.db.GetAuthor(ctx, authorID)
	if err != nil {
		s.logger.Error("GetAuthor sql call failed", zap.Error(err))
		return
	}
	return toAuthorProto(result)

}

func (s *Service) GetBook(ctx context.Context, in *pb.GetBookParams) (out *pb.Book, err error) {
	bookID := in.GetBookID()

	result, err := s.db.GetBook(ctx, bookID)
	if err != nil {
		s.logger.Error("GetBook sql call failed", zap.Error(err))
		return
	}
	return toBookProto(result)

}

func (s *Service) UpdateBook(ctx context.Context, in *pb.UpdateBookParams) (out *emptypb.Empty, err error) {
	arg, err := toUpdateBookParams(in)
	if err != nil {
		s.logger.Error("UpdateBook input adapter failed", zap.Error(err))
		return
	}

	err = s.db.UpdateBook(ctx, arg)
	if err != nil {
		s.logger.Error("UpdateBook sql call failed", zap.Error(err))
		return
	}
	return &emptypb.Empty{}, nil

}

func (s *Service) UpdateBookISBN(ctx context.Context, in *pb.UpdateBookISBNParams) (out *emptypb.Empty, err error) {
	arg, err := toUpdateBookISBNParams(in)
	if err != nil {
		s.logger.Error("UpdateBookISBN input adapter failed", zap.Error(err))
		return
	}

	err = s.db.UpdateBookISBN(ctx, arg)
	if err != nil {
		s.logger.Error("UpdateBookISBN sql call failed", zap.Error(err))
		return
	}
	return &emptypb.Empty{}, nil

}
