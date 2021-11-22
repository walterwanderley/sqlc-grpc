package books

import (
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"booktest/internal/validation"
	pb "booktest/proto/books"
)

func toBooksByTitleYearParams(in *pb.BooksByTitleYearParams) (out BooksByTitleYearParams, err error) {
	if in == nil {
		return
	}
	out.Title = in.GetTitle()
	out.Year = in.GetYear()
	return
}

func toCreateBookParams(in *pb.CreateBookParams) (out CreateBookParams, err error) {
	if in == nil {
		return
	}
	out.AuthorID = in.GetAuthorID()
	out.Isbn = in.GetIsbn()
	out.BookType = BookType(in.GetBookType())
	out.Title = in.GetTitle()
	out.Year = in.GetYear()
	if v := in.GetAvailable(); v != nil {
		if err = v.CheckValid(); err != nil {
			err = fmt.Errorf("invalid Available: %s%w", err.Error(), validation.ErrUserInput)
			return
		}
		out.Available = v.AsTime()
	} else {
		err = fmt.Errorf("Available is required%w", validation.ErrUserInput)
		return
	}
	out.Tags = in.GetTags()
	return
}

func toUpdateBookISBNParams(in *pb.UpdateBookISBNParams) (out UpdateBookISBNParams, err error) {
	if in == nil {
		return
	}
	out.Title = in.GetTitle()
	out.Tags = in.GetTags()
	out.BookID = in.GetBookID()
	out.Isbn = in.GetIsbn()
	return
}

func toUpdateBookParams(in *pb.UpdateBookParams) (out UpdateBookParams, err error) {
	if in == nil {
		return
	}
	out.Title = in.GetTitle()
	out.Tags = in.GetTags()
	out.BookType = BookType(in.GetBookType())
	out.BookID = in.GetBookID()
	return
}

func toAuthorProto(in Author) (out *pb.Author, err error) {
	out = new(pb.Author)
	out.AuthorID = in.AuthorID
	out.Name = in.Name
	return
}

func toBookProto(in Book) (out *pb.Book, err error) {
	out = new(pb.Book)
	out.BookID = in.BookID
	out.AuthorID = in.AuthorID
	out.Isbn = in.Isbn
	out.BookType = string(in.BookType)
	out.Title = in.Title
	out.Year = in.Year
	out.Available = timestamppb.New(in.Available)
	out.Tags = in.Tags
	return
}

func toBooksByTagsRowProto(in BooksByTagsRow) (out *pb.BooksByTagsRow, err error) {
	out = new(pb.BooksByTagsRow)
	out.BookID = in.BookID
	out.Title = in.Title
	if in.Name.Valid {
		out.Name = wrapperspb.String(in.Name.String)
	}
	out.Isbn = in.Isbn
	out.Tags = in.Tags
	return
}
