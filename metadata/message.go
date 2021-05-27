package metadata

import (
	"fmt"
	"go/ast"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Message struct {
	Name      string
	AttrNames []string
	AttrTypes []string
}

func (m *Message) ProtoAttributes() string {
	var s strings.Builder
	for i, name := range m.AttrNames {
		s.WriteString(fmt.Sprintf("    %s %s = %d;\n", toProtoType(m.AttrTypes[i]), name, i+1))
	}
	return s.String()
}

func createMessage(name string, s *ast.StructType) *Message {
	names := make([]string, 0)
	types := make([]string, 0)
	for _, f := range s.Fields.List {
		types = append(types, exprToStr(f.Type))
		var name string
		if len(f.Names) > 0 {
			name = f.Names[0].Name
		}
		names = append(names, name)
	}
	return &Message{
		Name:      name,
		AttrNames: names,
		AttrTypes: types,
	}
}

func customType(typ string) bool {
	ru, _ := utf8.DecodeRuneInString(typ[0:1])
	return unicode.IsUpper(ru)
}
