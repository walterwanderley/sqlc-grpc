package application

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"{{ .GoModule}}/api"
	database "{{ .GoModule}}/{{.SrcPath}}"
	"{{.GoModule}}/internal/server"
)
	
type Service struct {
    api.Unimplemented{{ .Package | UpperFirst}}Server
	db *database.Queries
}
	
func NewService(db *database.Queries) *Service {
	return &Service{db: db}
}

{{ range .Services }}
func (s *Service) {{.Name}}(ctx context.Context, in *{{.MethodInputType}}) (out *{{.MethodOutputType}}, err error) {
	{{ range .InputGrpc}}{{ .}}
	{{end}}
	{{ .ReturnCallDatabase}} err {{if not .EmptyOutput}}:{{end}}= s.db.{{ .Name}}(ctx{{ .ParamsCallDatabase}})
	if err != nil {			
		return
	}

	out = new({{.MethodOutputType}})
	{{ range .OutputGrpc}}{{ .}}		
	{{end}}				
	return
}
{{ end }}
