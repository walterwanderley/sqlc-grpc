package books

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"booktest/internal/validation"
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
		return
	}

	out = new(pb.BooksByTagsResponse)
	for _, r := range result {
		var item pb.BooksByTagsRow
		item.BookID = r.BookID
		item.Title = r.Title
		item.Name = r.Name
		item.Isbn = r.Isbn
		item.Tags = r.Tags
		out.Value = append(out.Value, &item)
	}

	return
}

func (s *Service) BooksByTitleYear(ctx context.Context, in *pb.BooksByTitleYearParams) (out *pb.BooksByTitleYearResponse, err error) {
	var arg BooksByTitleYearParams
	arg.Title = in.GetTitle()
	arg.Year = in.GetYear()

	result, err := s.db.BooksByTitleYear(ctx, arg)
	if err != nil {
		return
	}

	out = new(pb.BooksByTitleYearResponse)
	for _, r := range result {
		var item pb.Book
		item.BookID = r.BookID
		item.AuthorID = r.AuthorID
		item.Isbn = r.Isbn
		item.BookType = string(r.BookType)
		item.Title = r.Title
		item.Year = r.Year
		item.Available = timestamppb.New(r.Available)
		item.Tags = r.Tags
		out.Value = append(out.Value, &item)
	}

	return
}

func (s *Service) CreateAuthor(ctx context.Context, in *pb.CreateAuthorParams) (out *pb.Author, err error) {
	name := in.GetName()

	result, err := s.db.CreateAuthor(ctx, name)
	if err != nil {
		return
	}

	out = new(pb.Author)
	out.AuthorID = result.AuthorID
	out.Name = result.Name

	return
}

func (s *Service) CreateBook(ctx context.Context, in *pb.CreateBookParams) (out *pb.Book, err error) {
	var arg CreateBookParams
	arg.AuthorID = in.GetAuthorID()
	arg.Isbn = in.GetIsbn()
	arg.BookType = BookType(in.GetBookType())
	arg.Title = in.GetTitle()
	arg.Year = in.GetYear()
	if v := in.GetAvailable(); v != nil {
		if err = v.CheckValid(); err != nil {
			err = fmt.Errorf("invalid Available: %s%w", err.Error(), validation.ErrUserInput)
			return
		}
		arg.Available = v.AsTime()
	} else {
		err = fmt.Errorf("Available is required%w", validation.ErrUserInput)
		return
	}
	arg.Tags = in.GetTags()

	result, err := s.db.CreateBook(ctx, arg)
	if err != nil {
		return
	}

	out = new(pb.Book)
	out.BookID = result.BookID
	out.AuthorID = result.AuthorID
	out.Isbn = result.Isbn
	out.BookType = string(result.BookType)
	out.Title = result.Title
	out.Year = result.Year
	out.Available = timestamppb.New(result.Available)
	out.Tags = result.Tags

	return
}

func (s *Service) DeleteBook(ctx context.Context, in *pb.DeleteBookParams) (out *emptypb.Empty, err error) {
	bookID := in.GetBookID()

	err = s.db.DeleteBook(ctx, bookID)
	if err != nil {
		return
	}

	out = new(emptypb.Empty)

	return
}

func (s *Service) GetAuthor(ctx context.Context, in *pb.GetAuthorParams) (out *pb.Author, err error) {
	authorID := in.GetAuthorID()

	result, err := s.db.GetAuthor(ctx, authorID)
	if err != nil {
		return
	}

	out = new(pb.Author)
	out.AuthorID = result.AuthorID
	out.Name = result.Name

	return
}

func (s *Service) GetBook(ctx context.Context, in *pb.GetBookParams) (out *pb.Book, err error) {
	bookID := in.GetBookID()

	result, err := s.db.GetBook(ctx, bookID)
	if err != nil {
		return
	}

	out = new(pb.Book)
	out.BookID = result.BookID
	out.AuthorID = result.AuthorID
	out.Isbn = result.Isbn
	out.BookType = string(result.BookType)
	out.Title = result.Title
	out.Year = result.Year
	out.Available = timestamppb.New(result.Available)
	out.Tags = result.Tags

	return
}

func (s *Service) UpdateBook(ctx context.Context, in *pb.UpdateBookParams) (out *emptypb.Empty, err error) {
	var arg UpdateBookParams
	arg.Title = in.GetTitle()
	arg.Tags = in.GetTags()
	arg.BookType = BookType(in.GetBookType())
	arg.BookID = in.GetBookID()

	err = s.db.UpdateBook(ctx, arg)
	if err != nil {
		return
	}

	out = new(emptypb.Empty)

	return
}

func (s *Service) UpdateBookISBN(ctx context.Context, in *pb.UpdateBookISBNParams) (out *emptypb.Empty, err error) {
	var arg UpdateBookISBNParams
	arg.Title = in.GetTitle()
	arg.Tags = in.GetTags()
	arg.BookID = in.GetBookID()
	arg.Isbn = in.GetIsbn()

	err = s.db.UpdateBookISBN(ctx, arg)
	if err != nil {
		return
	}

	out = new(emptypb.Empty)

	return
}
