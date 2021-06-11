package metadata

import (
	"fmt"
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

func (m *Message) HasComplexAttribute() bool {
	for _, t := range m.AttrTypes {
		if customType(t) || strings.HasPrefix(t, "[]") {
			return true
		}
	}

	return false
}

func customType(typ string) bool {
	ru, _ := utf8.DecodeRuneInString(typ[0:1])
	return unicode.IsUpper(ru)
}
