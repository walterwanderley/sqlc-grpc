package metadata

import (
	"fmt"
	"strings"
)

type Service struct {
	Name       string
	InputNames []string
	InputTypes []string
	Output     []string
	Sql        string
	Messages   map[string]*Message
}

func (s *Service) MethodInputType() string {
	return fmt.Sprintf("pb.%sRequest", s.Name)
}

func (s *Service) MethodOutputType() string {
	switch {
	case s.HasCustomOutput():
		return fmt.Sprintf("pb.%s", canonicalName(s.Output[0]))
	default:
		return fmt.Sprintf("pb.%sResponse", s.Name)
	}
}

func (s *Service) ReturnCallDatabase() string {
	if !s.EmptyOutput() {
		return "result,"
	}
	return ""
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
		m := s.Messages[typ]
		for i, name := range m.AttrNames {
			attrName := UpperFirstCharacter(name)
			res = append(res, bindToGo("in", fmt.Sprintf("%s.%s", in, attrName), attrName, m.AttrTypes[i], false)...)
		}
	} else {
		for i, n := range s.InputNames {
			res = append(res, bindToGo("in", n, UpperFirstCharacter(n), s.InputTypes[i], true)...)
		}
	}

	return res
}

func (s *Service) OutputGrpc() []string {
	res := make([]string, 0)
	if s.HasArrayOutput() {
		res = append(res, fmt.Sprintf("out = new(%s)", s.MethodOutputType()))
		res = append(res, "for _, r := range result {")
		typ := canonicalName(s.Output[0])
		res = append(res, fmt.Sprintf("var item *pb.%s", typ))
		res = append(res, fmt.Sprintf("item, err = to%s(r)", typ))
		res = append(res, "if err != nil { return }")
		res = append(res, "out.Value = append(out.Value, item)")
		res = append(res, "}")
		res = append(res, "return")
		return res
	}

	if s.HasCustomOutput() {
		for _, n := range s.Output {
			res = append(res, fmt.Sprintf("return to%s(result)", canonicalName(n)))
		}
		return res
	}
	res = append(res, fmt.Sprintf("out = new(%s)", s.MethodOutputType()))
	if !s.EmptyOutput() {
		res = append(res, "out.Value = result")
	}
	res = append(res, "return")

	return res
}

func (s *Service) RpcSignature() string {
	var b strings.Builder
	b.WriteString(s.Name)
	b.WriteString("(")
	b.WriteString(fmt.Sprintf("%sRequest", s.Name))
	b.WriteString(") returns (")
	switch {
	case s.HasCustomOutput():
		b.WriteString(canonicalName(s.Output[0]))
	default:
		b.WriteString(fmt.Sprintf("%sResponse", s.Name))
	}
	b.WriteString(")")
	return b.String()
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

	if msg, ok := s.Messages[s.InputTypes[0]]; ok {
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

	return customType(s.Output[0])
}

func (s *Service) HasArrayOutput() bool {
	if s.EmptyOutput() {
		return false
	}
	return strings.HasPrefix(s.Output[0], "[]") && s.Output[0] != "[]byte"
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
	return len(s.Output) == 0
}

func (s *Service) ProtoOutputs() string {
	var b strings.Builder
	for i, name := range s.Output {
		fmt.Fprintf(&b, "    %s value = %d;\n", toProtoType(name), i+1)
	}
	return b.String()
}
