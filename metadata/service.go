package metadata

import (
	"fmt"
	"go/ast"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Service struct {
	Name       string
	InputNames []string
	InputTypes []string
	Output     []string
	Messages   map[string]*Message
}

func (s *Service) MethodInputType() string {
	switch {
	case s.EmptyInput():
		return "emptypb.Empty"
	default:
		return fmt.Sprintf("api.%sParams", s.Name)
	}
}

func (s *Service) MethodOutputType() string {
	switch {
	case s.EmptyOutput():
		return "emptypb.Empty"
	case s.HasCustomOutput():
		return fmt.Sprintf("api.%s", s.Output[0])
	default:
		return fmt.Sprintf("api.%sResponse", s.Name)
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
		res = append(res, fmt.Sprintf("var %s database.%s", in, typ))
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
	if s.EmptyOutput() {
		return res
	}

	if s.HasArrayOutput() {
		res = append(res, "for _, r := range result {")
		typ := strings.TrimPrefix(s.Output[0], "[]")
		res = append(res, fmt.Sprintf("var item api.%s", typ))
		m := s.Messages[typ]
		for i, attr := range m.AttrNames {
			res = append(res, bindToProto("r", "item", UpperFirstCharacter(attr), m.AttrTypes[i])...)
		}
		res = append(res, "out.Value = append(out.Value, &item)")
		res = append(res, "}")
		return res
	}

	if s.HasCustomOutput() {
		for _, n := range s.Output {
			m := s.Messages[n]
			for i, attr := range m.AttrNames {
				res = append(res, bindToProto("result", "out", UpperFirstCharacter(attr), m.AttrTypes[i])...)
			}
		}
		return res
	}
	if !s.EmptyOutput() {
		res = append(res, "out.Value = result")
		return res
	}

	return res
}

func bindToProto(src, dst, attrName, attrType string) []string {
	res := make([]string, 0)
	switch attrType {
	case "sql.NullBool":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = wrapperspb.Bool(%s.%s.Bool) }", dst, attrName, src, attrName))
	case "sql.NullInt32":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = wrapperspb.Int32(%s.%s.Int32) }", dst, attrName, src, attrName))
	case "sql.NullInt64":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = wrapperspb.Int64(%s.%s.Int64) }", dst, attrName, src, attrName))
	case "sql.NullFloat64":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = wrapperspb.Float64(%s.%s.Float64) }", dst, attrName, src, attrName))
	case "sql.NullString":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = wrapperspb.String(%s.%s.String) }", dst, attrName, src, attrName))
	case "sql.NullTime":
		res = append(res, fmt.Sprintf("if %s.%s.Valid {", src, attrName))
		res = append(res, fmt.Sprintf("%s.%s = timestamppb.New(%s.%s.Time) }", dst, attrName, src, attrName))
	case "time.Time":
		res = append(res, fmt.Sprintf("%s.%s = timestamppb.New(%s.%s)", dst, attrName, src, attrName))
	case "uuid.UUID":
		res = append(res, fmt.Sprintf("%s.%s = %s.%s.String()", dst, attrName, src, attrName))
	default:
		res = append(res, fmt.Sprintf("%s.%s = %s.%s", dst, attrName, src, attrName))
	}
	return res
}

func bindToGo(src, dst, attrName, attrType string, newVar bool) []string {
	res := make([]string, 0)
	switch attrType {
	case "sql.NullBool":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, attrName))
		res = append(res, fmt.Sprintf("%s = sql.NullBool{Valid: true, Bool: v.Value}", dst))
		res = append(res, "}")
	case "sql.NullInt32":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, attrName))
		res = append(res, fmt.Sprintf("%s = sql.NullInt32{Valid: true, Int32: v.Value}", dst))
		res = append(res, "}")
	case "sql.NullInt64":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, attrName))
		res = append(res, fmt.Sprintf("%s = sql.NullInt64{Valid: true, Int64: v.Value}", dst))
		res = append(res, "}")
	case "sql.NullFloat64":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, attrName))
		res = append(res, fmt.Sprintf("%s = sql.NullFloat64{Valid: true, Float64: v.Value}", dst))
		res = append(res, "}")
	case "sql.NullString":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, attrName))
		res = append(res, fmt.Sprintf("%s = sql.NullString{Valid: true, String: v.Value}", dst))
		res = append(res, "}")
	case "sql.NullTime":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, attrName))
		res = append(res, fmt.Sprintf("if err = v.CheckValid(); err != nil { err = fmt.Errorf(\"invalid %s: %%s%%w\", err.Error(), server.ErrUserInput)", attrName))
		res = append(res, "return }")
		res = append(res, "t := v.AsTime()")
		res = append(res, "if !t.IsZero() {")
		res = append(res, fmt.Sprintf("%s.Valid = true", dst))
		res = append(res, fmt.Sprintf("%s.Time = t } }", dst))
	case "time.Time":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if v := %s.Get%s(); v != nil {", src, attrName))
		res = append(res, fmt.Sprintf("if err = v.CheckValid(); err != nil { err = fmt.Errorf(\"invalid %s: %%s%%w\", err.Error(), server.ErrUserInput)", attrName))
		res = append(res, "return }")
		res = append(res, fmt.Sprintf("%s = v.AsTime()", dst))
		res = append(res, fmt.Sprintf("} else { err = fmt.Errorf(\"%s is required%%w\", server.ErrUserInput)", attrName))
		res = append(res, "return }")
	case "uuid.UUID":
		if newVar {
			res = append(res, fmt.Sprintf("var %s %s", dst, attrType))
		}
		res = append(res, fmt.Sprintf("if %s, err = uuid.Parse(%s.Get%s()); err != nil {", dst, src, attrName))
		res = append(res, fmt.Sprintf("err = fmt.Errorf(\"invalid %s: %%s%%w\", err.Error(), server.ErrUserInput)", attrName))
		res = append(res, "return }")
	default:
		if newVar {
			res = append(res, fmt.Sprintf("%s := %s.Get%s()", dst, src, attrName))
		} else {
			res = append(res, fmt.Sprintf("%s = %s.Get%s()", dst, src, attrName))
		}
	}
	return res
}

func (s *Service) RpcSignature() string {
	var b strings.Builder
	b.WriteString(s.Name)
	b.WriteString("(")
	switch {
	case s.EmptyInput():
		b.WriteString("google.protobuf.Empty")
	default:
		b.WriteString(fmt.Sprintf("%sParams", s.Name))
	}
	b.WriteString(") returns (")
	switch {
	case s.EmptyOutput():
		b.WriteString("google.protobuf.Empty")
	case s.HasCustomOutput():
		b.WriteString(s.Output[0])
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
	ru, _ := utf8.DecodeRuneInString(s.InputTypes[0][0:1])
	return unicode.IsUpper(ru)
}

func (s *Service) HasArrayParams() bool {
	if s.EmptyInput() {
		return false
	}

	return strings.HasPrefix(s.InputTypes[0], "[]")
}

func (s *Service) HasCustomOutput() bool {
	if s.EmptyOutput() {
		return false
	}
	ru, _ := utf8.DecodeRuneInString(s.Output[0][0:1])
	return unicode.IsUpper(ru)
}

func (s *Service) HasArrayOutput() bool {
	if s.EmptyOutput() {
		return false
	}
	return strings.HasPrefix(s.Output[0], "[]")
}

func (s *Service) ProtoInputs() string {
	var b strings.Builder
	for i, name := range s.InputNames {
		fmt.Fprintf(&b, "\n    %s %s = %d;", toProtoType(s.InputTypes[i]), name, i+1)
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

func visitFunc(fun *ast.FuncDecl, def *Definition) {
	if !isMethodValid(fun) {
		return
	}

	inputNames := make([]string, 0)
	inputTypes := make([]string, 0)
	output := make([]string, 0)

	// first param is always a context
	for i := 1; i < len(fun.Type.Params.List); i++ {
		p := fun.Type.Params.List[i]
		inputNames = append(inputNames, p.Names[0].Name)
		inputTypes = append(inputTypes, exprToStr(p.Type))
	}

	// last output result is always an error
	for i := 0; i < len(fun.Type.Results.List)-1; i++ {
		p := fun.Type.Results.List[0]
		output = append(output, exprToStr(p.Type))
	}

	def.Services = append(def.Services, &Service{
		Name:       fun.Name.String(),
		InputNames: inputNames,
		InputTypes: inputTypes,
		Output:     output,
		Messages:   def.Messages,
	})
}
func isMethodValid(fun *ast.FuncDecl) bool {
	if !fun.Name.IsExported() {
		return false
	}

	if fun.Recv == nil || len(fun.Recv.List) != 1 {
		return false
	}

	typ, ok := fun.Recv.List[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}

	if fun.Type.Params == nil || len(fun.Type.Params.List) == 0 ||
		fun.Type.Results == nil || len(fun.Type.Results.List) == 0 {
		return false
	}

	if exprToStr(fun.Type.Params.List[0].Type) != "context.Context" {
		return false
	}

	if exprToStr(fun.Type.Results.List[len(fun.Type.Results.List)-1].Type) != "error" {
		return false
	}

	t, ok := typ.X.(*ast.Ident)
	if !ok {
		return false
	}

	if t.Name != "Queries" {
		return false
	}
	return true
}
