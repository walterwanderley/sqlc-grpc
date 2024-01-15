package templates

import (
	"embed"
	"html/template"

	"github.com/walterwanderley/sqlc-grpc/converter"
	"github.com/walterwanderley/sqlc-grpc/metadata"
)

//go:embed *
var Files embed.FS

var Funcs = template.FuncMap{
	"PascalCase":          converter.ToPascalCase,
	"SnakeCase":           converter.ToSnakeCase,
	"UpperFirstCharacter": converter.UpperFirstCharacter,
	"Input":               metadata.InputGrpc,
	"Output":              metadata.OutputGrpc,
}
