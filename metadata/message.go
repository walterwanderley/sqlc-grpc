package metadata

import (
	"fmt"
	"go/ast"
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
		s.WriteString(fmt.Sprintf("    %s %s = %d;\n", toProtoType(m.AttrTypes[i]), lowerFirstCharacter(name), i+1))
	}
	return s.String()
}

func (m *Message) adjustType(messages map[string]*Message) {
	for i, t := range m.AttrTypes {
		m.AttrTypes[i] = adjustType(t, messages)
	}
}

func createStructMessage(name string, s *ast.StructType) *Message {
	names := make([]string, 0)
	types := make([]string, 0)
	for _, f := range s.Fields.List {
		if len(f.Names) == 0 || !firstIsUpper(f.Names[0].Name) {
			continue
		}
		types = append(types, exprToStr(f.Type))
		names = append(names, f.Names[0].Name)
	}
	return &Message{
		Name:      name,
		AttrNames: names,
		AttrTypes: types,
	}
}

func createArrayMessage(name string, s *ast.ArrayType) *Message {
	return &Message{
		Name:        name,
		IsArray:     true,
		ElementType: exprToStr(s.Elt),
	}
}

func createAliasMessage(name string, s *ast.Ident) *Message {
	return &Message{
		Name:        name,
		ElementType: exprToStr(s),
	}
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
