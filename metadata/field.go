package metadata

import (
	"fmt"
	"strings"
)

type Field struct {
	Name                string
	Type                string
	CustomProtoComments []string
	CustomProtoOptions  []string
}

func (f *Field) Proto(tag int) string {
	var sb strings.Builder
	for _, line := range f.CustomProtoComments {
		sb.WriteString(fmt.Sprintf("    // %s\n", line))
	}
	sb.WriteString(fmt.Sprintf("    %s %s = %d%s;\n", ToProtoType(f.Type), ToSnakeCase(f.Name), tag, f.formatProtoOptions()))
	return sb.String()
}

func (f *Field) formatProtoOptions() string {
	var sb strings.Builder
	if len(f.CustomProtoOptions) > 0 {
		sb.WriteString(" [")
		for i, opt := range f.CustomProtoOptions {
			if i > 1 {
				sb.WriteString("\n")
			}
			sb.WriteString(opt)
		}
		sb.WriteString("]")
	}
	return sb.String()
}
