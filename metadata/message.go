package metadata

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Message struct {
	Name        string
	AttrNames   []string
	AttrTypes   []string
	IsArray     bool
	ElementType string
}

func (m *Message) ProtoAttributes() string {
	var s strings.Builder
	for i, name := range m.AttrNames {
		s.WriteString(fmt.Sprintf("    %s %s = %d;\n", toProtoType(m.AttrTypes[i]), ToSnakeCase(name), i+1))
	}
	return s.String()
}

func (m *Message) adjustType(messages map[string]*Message) {
	for i, t := range m.AttrTypes {
		m.AttrTypes[i] = adjustType(t, messages)
	}
}

func createStructMessage(name string, s *ast.StructType) (*Message, error) {
	names := make([]string, 0)
	types := make([]string, 0)
	for _, f := range s.Fields.List {
		if len(f.Names) == 0 || !firstIsUpper(f.Names[0].Name) {
			continue
		}
		typ, err := exprToStr(f.Type)
		if err != nil {
			return nil, err
		}
		types = append(types, typ)
		names = append(names, f.Names[0].Name)
	}
	return &Message{
		Name:      name,
		AttrNames: names,
		AttrTypes: types,
	}, nil
}

func createArrayMessage(name string, s *ast.ArrayType) (*Message, error) {
	elt, err := exprToStr(s.Elt)
	if err != nil {
		return nil, err
	}
	return &Message{
		Name:        name,
		IsArray:     true,
		ElementType: elt,
	}, nil
}

func createAliasMessage(name string, s *ast.Ident) (*Message, error) {
	str, err := exprToStr(s)
	if err != nil {
		return nil, err
	}
	return &Message{
		Name:        name,
		ElementType: str,
	}, nil
}

func customType(typ string) bool {
	typ = strings.TrimPrefix(typ, "*")
	return firstIsUpper(typ)
}

func firstIsUpper(s string) bool {
	ru, _ := utf8.DecodeRuneInString(s[0:1])
	return unicode.IsUpper(ru)
}

func adjustType(typ string, messages map[string]*Message) string {
	if m, ok := messages[typ]; ok {
		var prefix string
		if m.IsArray {
			prefix = "[]"
		}
		if m.ElementType != "" {
			return prefix + typ + "." + m.ElementType
		}
	}

	return typ
}

func OriginalAndElementType(typ string) (original, element string) {
	typ = strings.TrimPrefix(typ, "[]")
	t := strings.Split(typ, ".")
	return t[0], strings.Join(t[1:], ".")
}

func (m *Message) HasComplexAttribute() bool {
	for _, t := range m.AttrTypes {
		if customType(t) || strings.HasPrefix(t, "[]") {
			return true
		}
	}

	return false
}

func (m *Message) AdapterToGo(src, dst string) []string {
	res := make([]string, 0)
	for i, attr := range m.AttrNames {
		attrName := UpperFirstCharacter(attr)
		res = append(res, bindToGo(src, fmt.Sprintf("%s.%s", dst, attrName), attrName, m.AttrTypes[i], false)...)
	}
	return res
}

func (m *Message) AdapterToProto(src, dst string) []string {
	res := make([]string, 0)
	for i, attr := range m.AttrNames {
		res = append(res, bindToProto(src, dst, UpperFirstCharacter(attr), m.AttrTypes[i])...)
	}
	return res
}

func (m *Message) ProtoName() string {
	return regexp.MustCompile("Params$").ReplaceAllString(m.Name, "Request")
}
