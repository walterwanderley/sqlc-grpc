package metadata

import (
	"fmt"
	"strings"

	"github.com/walterwanderley/sqlc-grpc/converter"
)

func InputGrpc(s *Service) []string {
	res := make([]string, 0)
	if s.EmptyInput() {
		return res
	}

	if s.HasCustomParams() {
		typ := s.InputTypes[0]
		in := s.InputNames[0]
		if strings.HasPrefix(typ, "*") {
			res = append(res, fmt.Sprintf("%s := new(%s)", in, typ[1:]))
		} else {
			res = append(res, fmt.Sprintf("var %s %s", in, typ))
		}
		m := s.Messages[converter.CanonicalName(typ)]
		for _, f := range m.Fields {
			attrName := converter.UpperFirstCharacter(f.Name)
			res = append(res, converter.BindToGo("req", fmt.Sprintf("%s.%s", in, attrName), attrName, f.Type, false)...)
		}
	} else {
		for i, n := range s.InputNames {
			res = append(res, converter.BindToGo("req", n, converter.UpperFirstCharacter(n), s.InputTypes[i], true)...)
		}
	}

	return res
}

func OutputGrpc(s *Service) []string {
	name := converter.UpperFirstCharacter(s.Name)
	res := make([]string, 0)
	if s.HasArrayOutput() {
		res = append(res, fmt.Sprintf("res := new(pb.%sResponse)", name))
		res = append(res, "for _, r := range result {")
		res = append(res, fmt.Sprintf("res.List = append(res.List, to%s(r))", converter.CanonicalName(s.Output)))
		res = append(res, "}")
		res = append(res, "return res, nil")
		return res
	}

	if s.HasCustomOutput() {
		res = append(res, fmt.Sprintf("return &pb.%sResponse{%s: to%s(result)}, nil", name, converter.CamelCaseProto(converter.CanonicalName(s.Output)), converter.CanonicalName(s.Output)))
		return res
	}
	if s.EmptyOutput() {
		res = append(res, fmt.Sprintf("return &pb.%sResponse{}, nil", name))
		return res
	} else {
		if s.Output == "sql.Result" {
			res = append(res, fmt.Sprintf("return &pb.%sResponse{Value: toExecResult(result)}, nil", name))
			return res
		}
		if s.Output == "uuid.UUID" {
			res = append(res, fmt.Sprintf("return &pb.%sResponse{Value: result.String()}, nil", name))
			return res
		}
		if s.Output == "pgtype.UUID" {
			res = append(res, "if uuidStr, err := result.MarshalJSON(); err != nil {")
			res = append(res, "return nil, fmt.Errorf(\"failed to convert UUID to string: %w\", err)")
			res = append(res, "} else {")
			res = append(res, fmt.Sprintf("return &pb.%sResponse{Value: wrapperspb.String(string(uuidStr))}, nil", name))
			res = append(res, "}")
			return res
		}
		res = append(res, fmt.Sprintf("return &pb.%sResponse{Value: result}, nil", name))
	}

	return res
}
