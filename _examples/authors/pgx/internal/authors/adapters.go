// Code generated by sqlc-grpc (https://github.com/walterwanderley/sqlc-grpc). DO NOT EDIT.

package authors

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pb "authors/api/authors/v1"
)

func toAuthors(in *Authors) *pb.Authors {
	if in == nil {
		return nil
	}
	out := new(pb.Authors)
	out.Id = in.ID
	out.Name = in.Name
	if in.Bio.Valid {
		out.Bio = wrapperspb.String(in.Bio.String)
	}
	if in.CreatedAt.Valid {
		out.CreatedAt = timestamppb.New(in.CreatedAt.Time)
	}
	return out
}
