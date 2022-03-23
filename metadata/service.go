package metadata

import (
	"fmt"
	"strings"
)

type Service struct {
	Name              string
	InputNames        []string
	InputTypes        []string
	Output            string
	Sql               string
	Messages          map[string]*Message
	CustomHttpOptions []string
}

func (s *Service) ParamsCallDatabase() string {
	if s.EmptyInput() {
		return ""
	}
	return ", " + strings.Join(s.InputNames, ", ")
}

func (s *Service) InputGrpc() []string {
	res := make([]string, 0)
	if s.EmptyInput() {
		return res
	}

	if s.HasCustomParams() {
		typ := s.InputTypes[0]
		in := s.InputNames[0]
		res = append(res, fmt.Sprintf("var %s %s", in, typ))
		m := s.Messages[canonicalName(typ)]
		for i, name := range m.AttrNames {
			attrName := UpperFirstCharacter(name)
			res = append(res, bindToGo("req", fmt.Sprintf("%s.%s", in, attrName), attrName, m.AttrTypes[i], false)...)
		}
	} else {
		for i, n := range s.InputNames {
			res = append(res, bindToGo("req", n, UpperFirstCharacter(n), s.InputTypes[i], true)...)
		}
	}

	return res
}

func (s *Service) OutputGrpc() []string {
	res := make([]string, 0)
	if s.HasArrayOutput() {
		res = append(res, fmt.Sprintf("res := new(pb.%sResponse)", s.Name))
		res = append(res, "for _, r := range result {")
		res = append(res, fmt.Sprintf("res.List = append(res.List, to%s(r))", canonicalName(s.Output)))
		res = append(res, "}")
		res = append(res, "return res, nil")
		return res
	}

	if s.HasCustomOutput() {
		res = append(res, fmt.Sprintf("return &pb.%sResponse{%s: to%s(result)}, nil", s.Name, camelCaseProto(canonicalName(s.Output)), canonicalName(s.Output)))
		return res
	}
	if s.EmptyOutput() {
		res = append(res, fmt.Sprintf("return &pb.%sResponse{}, nil", s.Name))
	} else {
		res = append(res, fmt.Sprintf("return &pb.%sResponse{Value: result}, nil", s.Name))
	}

	return res
}

func (s *Service) HasCustomParams() bool {
	if s.EmptyInput() {
		return false
	}

	return customType(s.InputTypes[0])
}

func (s *Service) HasSimpleParams() bool {
	if s.HasArrayParams() {
		return false
	}

	if !s.HasCustomParams() || s.EmptyInput() {
		return true
	}

	if msg, ok := s.Messages[canonicalName(s.InputTypes[0])]; ok {
		return !msg.HasComplexAttribute()
	}

	return false
}

func (s *Service) HasArrayParams() bool {
	if s.EmptyInput() {
		return false
	}

	return strings.HasPrefix(s.InputTypes[0], "[]") && s.InputTypes[0] != "[]byte"
}

func (s *Service) HasCustomOutput() bool {
	if s.EmptyOutput() {
		return false
	}

	return customType(s.Output)
}

func (s *Service) HasArrayOutput() bool {
	if s.EmptyOutput() {
		return false
	}
	return strings.HasPrefix(s.Output, "[]") && s.Output != "[]byte"
}

func (s *Service) ProtoInputs() string {
	var b strings.Builder
	for i, name := range s.InputNames {
		fmt.Fprintf(&b, "\n    %s %s = %d;", toProtoType(s.InputTypes[i]), ToSnakeCase(name), i+1)
	}
	return b.String()
}

func (s *Service) EmptyInput() bool {
	return len(s.InputTypes) == 0
}

func (s *Service) EmptyOutput() bool {
	return s.Output == ""
}

func (s *Service) ProtoOutputs() string {
	if s.EmptyOutput() {
		return ""
	}
	name := "value"
	if s.HasArrayOutput() {
		name = "list"
	} else if s.HasCustomOutput() {
		name = ToSnakeCase(canonicalName(s.Output))
	}
	return fmt.Sprintf("    %s %s = 1;\n", toProtoType(s.Output), name)
}
