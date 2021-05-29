package {{.Package}}

import (
	"context"
	"database/sql"
	"encoding/json"
	"net"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pb "{{ .GoModule}}/proto/{{.Package}}"
	"{{.GoModule}}/internal/validation"
)
	
type Service struct {
    pb.Unimplemented{{ .Package | UpperFirst}}Server
	db *Queries
}
	
func NewService(db *Queries) *Service {
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
