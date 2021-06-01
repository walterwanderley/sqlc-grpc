syntax = "proto3";

package {{.Package}};

option go_package = "{{.GoModule}}/proto/{{.Package}}";

{{range .ProtoImports}}{{ .}}
{{end}}
service {{.Package}} {
    {{range .Services}}
    rpc {{.RpcSignature}} { }
    {{end}}
}
{{range .Services}}{{if and (not .HasCustomParams) (not .EmptyInput)}}
message {{.Name}}Params { {{.ProtoInputs}}
}
{{end}}{{if and (not .HasCustomOutput) (not .EmptyOutput)}}
message {{.Name}}Response {
{{.ProtoOutputs -}}    
}{{end}}{{end}}
{{ range $key, $value := .Messages}}
message {{$key}} {
{{$value.ProtoAttributes -}}
}
{{end}}