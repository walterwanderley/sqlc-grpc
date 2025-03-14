// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: books/v1/books.proto

/*
Package v1 is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package v1

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Suppress "imported and not used" errors
var (
	_ codes.Code
	_ io.Reader
	_ status.Status
	_ = errors.New
	_ = runtime.String
	_ = utilities.NewDoubleArray
	_ = metadata.Join
)

func request_BooksService_BooksByTags_0(ctx context.Context, marshaler runtime.Marshaler, client BooksServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq BooksByTagsRequest
		metadata runtime.ServerMetadata
	)
	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq.Dollar_1); err != nil && !errors.Is(err, io.EOF) {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	msg, err := client.BooksByTags(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err
}

func local_request_BooksService_BooksByTags_0(ctx context.Context, marshaler runtime.Marshaler, server BooksServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq BooksByTagsRequest
		metadata runtime.ServerMetadata
	)
	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq.Dollar_1); err != nil && !errors.Is(err, io.EOF) {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	msg, err := server.BooksByTags(ctx, &protoReq)
	return msg, metadata, err
}

var filter_BooksService_BooksByTitleYear_0 = &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}

func request_BooksService_BooksByTitleYear_0(ctx context.Context, marshaler runtime.Marshaler, client BooksServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq BooksByTitleYearRequest
		metadata runtime.ServerMetadata
	)
	io.Copy(io.Discard, req.Body)
	if err := req.ParseForm(); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	if err := runtime.PopulateQueryParameters(&protoReq, req.Form, filter_BooksService_BooksByTitleYear_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	msg, err := client.BooksByTitleYear(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err
}

func local_request_BooksService_BooksByTitleYear_0(ctx context.Context, marshaler runtime.Marshaler, server BooksServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq BooksByTitleYearRequest
		metadata runtime.ServerMetadata
	)
	if err := req.ParseForm(); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	if err := runtime.PopulateQueryParameters(&protoReq, req.Form, filter_BooksService_BooksByTitleYear_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	msg, err := server.BooksByTitleYear(ctx, &protoReq)
	return msg, metadata, err
}

func request_BooksService_CreateAuthor_0(ctx context.Context, marshaler runtime.Marshaler, client BooksServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq CreateAuthorRequest
		metadata runtime.ServerMetadata
	)
	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && !errors.Is(err, io.EOF) {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	msg, err := client.CreateAuthor(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err
}

func local_request_BooksService_CreateAuthor_0(ctx context.Context, marshaler runtime.Marshaler, server BooksServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq CreateAuthorRequest
		metadata runtime.ServerMetadata
	)
	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && !errors.Is(err, io.EOF) {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	msg, err := server.CreateAuthor(ctx, &protoReq)
	return msg, metadata, err
}

func request_BooksService_CreateBook_0(ctx context.Context, marshaler runtime.Marshaler, client BooksServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq CreateBookRequest
		metadata runtime.ServerMetadata
	)
	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && !errors.Is(err, io.EOF) {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	msg, err := client.CreateBook(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err
}

func local_request_BooksService_CreateBook_0(ctx context.Context, marshaler runtime.Marshaler, server BooksServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq CreateBookRequest
		metadata runtime.ServerMetadata
	)
	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && !errors.Is(err, io.EOF) {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	msg, err := server.CreateBook(ctx, &protoReq)
	return msg, metadata, err
}

func request_BooksService_DeleteBook_0(ctx context.Context, marshaler runtime.Marshaler, client BooksServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq DeleteBookRequest
		metadata runtime.ServerMetadata
		err      error
	)
	io.Copy(io.Discard, req.Body)
	val, ok := pathParams["book_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "book_id")
	}
	protoReq.BookId, err = runtime.Int32(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "book_id", err)
	}
	msg, err := client.DeleteBook(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err
}

func local_request_BooksService_DeleteBook_0(ctx context.Context, marshaler runtime.Marshaler, server BooksServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq DeleteBookRequest
		metadata runtime.ServerMetadata
		err      error
	)
	val, ok := pathParams["book_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "book_id")
	}
	protoReq.BookId, err = runtime.Int32(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "book_id", err)
	}
	msg, err := server.DeleteBook(ctx, &protoReq)
	return msg, metadata, err
}

func request_BooksService_GetAuthor_0(ctx context.Context, marshaler runtime.Marshaler, client BooksServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq GetAuthorRequest
		metadata runtime.ServerMetadata
		err      error
	)
	io.Copy(io.Discard, req.Body)
	val, ok := pathParams["author_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "author_id")
	}
	protoReq.AuthorId, err = runtime.Int32(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "author_id", err)
	}
	msg, err := client.GetAuthor(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err
}

func local_request_BooksService_GetAuthor_0(ctx context.Context, marshaler runtime.Marshaler, server BooksServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq GetAuthorRequest
		metadata runtime.ServerMetadata
		err      error
	)
	val, ok := pathParams["author_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "author_id")
	}
	protoReq.AuthorId, err = runtime.Int32(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "author_id", err)
	}
	msg, err := server.GetAuthor(ctx, &protoReq)
	return msg, metadata, err
}

func request_BooksService_GetBook_0(ctx context.Context, marshaler runtime.Marshaler, client BooksServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq GetBookRequest
		metadata runtime.ServerMetadata
		err      error
	)
	io.Copy(io.Discard, req.Body)
	val, ok := pathParams["book_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "book_id")
	}
	protoReq.BookId, err = runtime.Int32(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "book_id", err)
	}
	msg, err := client.GetBook(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err
}

func local_request_BooksService_GetBook_0(ctx context.Context, marshaler runtime.Marshaler, server BooksServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq GetBookRequest
		metadata runtime.ServerMetadata
		err      error
	)
	val, ok := pathParams["book_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "book_id")
	}
	protoReq.BookId, err = runtime.Int32(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "book_id", err)
	}
	msg, err := server.GetBook(ctx, &protoReq)
	return msg, metadata, err
}

func request_BooksService_UpdateBook_0(ctx context.Context, marshaler runtime.Marshaler, client BooksServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq UpdateBookRequest
		metadata runtime.ServerMetadata
		err      error
	)
	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && !errors.Is(err, io.EOF) {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	val, ok := pathParams["book_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "book_id")
	}
	protoReq.BookId, err = runtime.Int32(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "book_id", err)
	}
	msg, err := client.UpdateBook(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err
}

func local_request_BooksService_UpdateBook_0(ctx context.Context, marshaler runtime.Marshaler, server BooksServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq UpdateBookRequest
		metadata runtime.ServerMetadata
		err      error
	)
	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && !errors.Is(err, io.EOF) {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	val, ok := pathParams["book_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "book_id")
	}
	protoReq.BookId, err = runtime.Int32(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "book_id", err)
	}
	msg, err := server.UpdateBook(ctx, &protoReq)
	return msg, metadata, err
}

func request_BooksService_UpdateBookISBN_0(ctx context.Context, marshaler runtime.Marshaler, client BooksServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq UpdateBookISBNRequest
		metadata runtime.ServerMetadata
		err      error
	)
	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && !errors.Is(err, io.EOF) {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	val, ok := pathParams["book_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "book_id")
	}
	protoReq.BookId, err = runtime.Int32(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "book_id", err)
	}
	msg, err := client.UpdateBookISBN(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err
}

func local_request_BooksService_UpdateBookISBN_0(ctx context.Context, marshaler runtime.Marshaler, server BooksServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var (
		protoReq UpdateBookISBNRequest
		metadata runtime.ServerMetadata
		err      error
	)
	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && !errors.Is(err, io.EOF) {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	val, ok := pathParams["book_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "book_id")
	}
	protoReq.BookId, err = runtime.Int32(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "book_id", err)
	}
	msg, err := server.UpdateBookISBN(ctx, &protoReq)
	return msg, metadata, err
}

// RegisterBooksServiceHandlerServer registers the http handlers for service BooksService to "mux".
// UnaryRPC     :call BooksServiceServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterBooksServiceHandlerFromEndpoint instead.
// GRPC interceptors will not work for this type of registration. To use interceptors, you must use the "runtime.WithMiddlewares" option in the "runtime.NewServeMux" call.
func RegisterBooksServiceHandlerServer(ctx context.Context, mux *runtime.ServeMux, server BooksServiceServer) error {
	mux.Handle(http.MethodPost, pattern_BooksService_BooksByTags_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/books.v1.BooksService/BooksByTags", runtime.WithHTTPPathPattern("/books-by-tags"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_BooksService_BooksByTags_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_BooksService_BooksByTags_0(annotatedContext, mux, outboundMarshaler, w, req, response_BooksService_BooksByTags_0{resp.(*BooksByTagsResponse)}, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodGet, pattern_BooksService_BooksByTitleYear_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/books.v1.BooksService/BooksByTitleYear", runtime.WithHTTPPathPattern("/books"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_BooksService_BooksByTitleYear_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_BooksService_BooksByTitleYear_0(annotatedContext, mux, outboundMarshaler, w, req, response_BooksService_BooksByTitleYear_0{resp.(*BooksByTitleYearResponse)}, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodPost, pattern_BooksService_CreateAuthor_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/books.v1.BooksService/CreateAuthor", runtime.WithHTTPPathPattern("/authors"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_BooksService_CreateAuthor_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_BooksService_CreateAuthor_0(annotatedContext, mux, outboundMarshaler, w, req, response_BooksService_CreateAuthor_0{resp.(*CreateAuthorResponse)}, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodPost, pattern_BooksService_CreateBook_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/books.v1.BooksService/CreateBook", runtime.WithHTTPPathPattern("/books"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_BooksService_CreateBook_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_BooksService_CreateBook_0(annotatedContext, mux, outboundMarshaler, w, req, response_BooksService_CreateBook_0{resp.(*CreateBookResponse)}, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodDelete, pattern_BooksService_DeleteBook_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/books.v1.BooksService/DeleteBook", runtime.WithHTTPPathPattern("/books/{book_id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_BooksService_DeleteBook_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_BooksService_DeleteBook_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodGet, pattern_BooksService_GetAuthor_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/books.v1.BooksService/GetAuthor", runtime.WithHTTPPathPattern("/authors/{author_id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_BooksService_GetAuthor_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_BooksService_GetAuthor_0(annotatedContext, mux, outboundMarshaler, w, req, response_BooksService_GetAuthor_0{resp.(*GetAuthorResponse)}, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodGet, pattern_BooksService_GetBook_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/books.v1.BooksService/GetBook", runtime.WithHTTPPathPattern("/books/{book_id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_BooksService_GetBook_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_BooksService_GetBook_0(annotatedContext, mux, outboundMarshaler, w, req, response_BooksService_GetBook_0{resp.(*GetBookResponse)}, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodPut, pattern_BooksService_UpdateBook_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/books.v1.BooksService/UpdateBook", runtime.WithHTTPPathPattern("/books/{book_id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_BooksService_UpdateBook_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_BooksService_UpdateBook_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodPatch, pattern_BooksService_UpdateBookISBN_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/books.v1.BooksService/UpdateBookISBN", runtime.WithHTTPPathPattern("/books/{book_id}/isbn"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_BooksService_UpdateBookISBN_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_BooksService_UpdateBookISBN_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})

	return nil
}

// RegisterBooksServiceHandlerFromEndpoint is same as RegisterBooksServiceHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterBooksServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.NewClient(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Errorf("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Errorf("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()
	return RegisterBooksServiceHandler(ctx, mux, conn)
}

// RegisterBooksServiceHandler registers the http handlers for service BooksService to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterBooksServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterBooksServiceHandlerClient(ctx, mux, NewBooksServiceClient(conn))
}

// RegisterBooksServiceHandlerClient registers the http handlers for service BooksService
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "BooksServiceClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "BooksServiceClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "BooksServiceClient" to call the correct interceptors. This client ignores the HTTP middlewares.
func RegisterBooksServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client BooksServiceClient) error {
	mux.Handle(http.MethodPost, pattern_BooksService_BooksByTags_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateContext(ctx, mux, req, "/books.v1.BooksService/BooksByTags", runtime.WithHTTPPathPattern("/books-by-tags"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_BooksService_BooksByTags_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_BooksService_BooksByTags_0(annotatedContext, mux, outboundMarshaler, w, req, response_BooksService_BooksByTags_0{resp.(*BooksByTagsResponse)}, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodGet, pattern_BooksService_BooksByTitleYear_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateContext(ctx, mux, req, "/books.v1.BooksService/BooksByTitleYear", runtime.WithHTTPPathPattern("/books"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_BooksService_BooksByTitleYear_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_BooksService_BooksByTitleYear_0(annotatedContext, mux, outboundMarshaler, w, req, response_BooksService_BooksByTitleYear_0{resp.(*BooksByTitleYearResponse)}, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodPost, pattern_BooksService_CreateAuthor_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateContext(ctx, mux, req, "/books.v1.BooksService/CreateAuthor", runtime.WithHTTPPathPattern("/authors"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_BooksService_CreateAuthor_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_BooksService_CreateAuthor_0(annotatedContext, mux, outboundMarshaler, w, req, response_BooksService_CreateAuthor_0{resp.(*CreateAuthorResponse)}, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodPost, pattern_BooksService_CreateBook_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateContext(ctx, mux, req, "/books.v1.BooksService/CreateBook", runtime.WithHTTPPathPattern("/books"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_BooksService_CreateBook_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_BooksService_CreateBook_0(annotatedContext, mux, outboundMarshaler, w, req, response_BooksService_CreateBook_0{resp.(*CreateBookResponse)}, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodDelete, pattern_BooksService_DeleteBook_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateContext(ctx, mux, req, "/books.v1.BooksService/DeleteBook", runtime.WithHTTPPathPattern("/books/{book_id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_BooksService_DeleteBook_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_BooksService_DeleteBook_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodGet, pattern_BooksService_GetAuthor_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateContext(ctx, mux, req, "/books.v1.BooksService/GetAuthor", runtime.WithHTTPPathPattern("/authors/{author_id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_BooksService_GetAuthor_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_BooksService_GetAuthor_0(annotatedContext, mux, outboundMarshaler, w, req, response_BooksService_GetAuthor_0{resp.(*GetAuthorResponse)}, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodGet, pattern_BooksService_GetBook_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateContext(ctx, mux, req, "/books.v1.BooksService/GetBook", runtime.WithHTTPPathPattern("/books/{book_id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_BooksService_GetBook_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_BooksService_GetBook_0(annotatedContext, mux, outboundMarshaler, w, req, response_BooksService_GetBook_0{resp.(*GetBookResponse)}, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodPut, pattern_BooksService_UpdateBook_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateContext(ctx, mux, req, "/books.v1.BooksService/UpdateBook", runtime.WithHTTPPathPattern("/books/{book_id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_BooksService_UpdateBook_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_BooksService_UpdateBook_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})
	mux.Handle(http.MethodPatch, pattern_BooksService_UpdateBookISBN_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		annotatedContext, err := runtime.AnnotateContext(ctx, mux, req, "/books.v1.BooksService/UpdateBookISBN", runtime.WithHTTPPathPattern("/books/{book_id}/isbn"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_BooksService_UpdateBookISBN_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}
		forward_BooksService_UpdateBookISBN_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
	})
	return nil
}

type response_BooksService_BooksByTags_0 struct {
	*BooksByTagsResponse
}

func (m response_BooksService_BooksByTags_0) XXX_ResponseBody() interface{} {
	return m.List
}

type response_BooksService_BooksByTitleYear_0 struct {
	*BooksByTitleYearResponse
}

func (m response_BooksService_BooksByTitleYear_0) XXX_ResponseBody() interface{} {
	return m.List
}

type response_BooksService_CreateAuthor_0 struct {
	*CreateAuthorResponse
}

func (m response_BooksService_CreateAuthor_0) XXX_ResponseBody() interface{} {
	return m.Author
}

type response_BooksService_CreateBook_0 struct {
	*CreateBookResponse
}

func (m response_BooksService_CreateBook_0) XXX_ResponseBody() interface{} {
	return m.Book
}

type response_BooksService_GetAuthor_0 struct {
	*GetAuthorResponse
}

func (m response_BooksService_GetAuthor_0) XXX_ResponseBody() interface{} {
	return m.Author
}

type response_BooksService_GetBook_0 struct {
	*GetBookResponse
}

func (m response_BooksService_GetBook_0) XXX_ResponseBody() interface{} {
	return m.Book
}

var (
	pattern_BooksService_BooksByTags_0      = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0}, []string{"books-by-tags"}, ""))
	pattern_BooksService_BooksByTitleYear_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0}, []string{"books"}, ""))
	pattern_BooksService_CreateAuthor_0     = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0}, []string{"authors"}, ""))
	pattern_BooksService_CreateBook_0       = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0}, []string{"books"}, ""))
	pattern_BooksService_DeleteBook_0       = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 1, 0, 4, 1, 5, 1}, []string{"books", "book_id"}, ""))
	pattern_BooksService_GetAuthor_0        = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 1, 0, 4, 1, 5, 1}, []string{"authors", "author_id"}, ""))
	pattern_BooksService_GetBook_0          = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 1, 0, 4, 1, 5, 1}, []string{"books", "book_id"}, ""))
	pattern_BooksService_UpdateBook_0       = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 1, 0, 4, 1, 5, 1}, []string{"books", "book_id"}, ""))
	pattern_BooksService_UpdateBookISBN_0   = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 1, 0, 4, 1, 5, 1, 2, 2}, []string{"books", "book_id", "isbn"}, ""))
)

var (
	forward_BooksService_BooksByTags_0      = runtime.ForwardResponseMessage
	forward_BooksService_BooksByTitleYear_0 = runtime.ForwardResponseMessage
	forward_BooksService_CreateAuthor_0     = runtime.ForwardResponseMessage
	forward_BooksService_CreateBook_0       = runtime.ForwardResponseMessage
	forward_BooksService_DeleteBook_0       = runtime.ForwardResponseMessage
	forward_BooksService_GetAuthor_0        = runtime.ForwardResponseMessage
	forward_BooksService_GetBook_0          = runtime.ForwardResponseMessage
	forward_BooksService_UpdateBook_0       = runtime.ForwardResponseMessage
	forward_BooksService_UpdateBookISBN_0   = runtime.ForwardResponseMessage
)
