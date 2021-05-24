package metadata

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type Definition struct {
	Database string
	Package  string
	GoModule string
	SrcPath  string
	Services []*Service
	Messages map[string]*Message
}

func (d *Definition) ProtoImports() []string {
	r := make([]string, 0)
	if d.importEmpty() {
		r = append(r, `import "google/protobuf/empty.proto";`)
	}
	if d.importTimestamp() {
		r = append(r, `import "google/protobuf/timestamp.proto";`)
	}
	if d.importWrappers() {
		r = append(r, `import "google/protobuf/wrappers.proto";`)
	}
	return r
}

func (d *Definition) importEmpty() bool {
	for _, s := range d.Services {
		if s.EmptyInput() || s.EmptyOutput() {
			return true
		}
	}
	return false
}

func (d *Definition) importTimestamp() bool {
	for _, m := range d.Messages {
		for _, typ := range m.AttrTypes {
			if typ == "time.Time" || typ == "sql.NullTime" {
				return true
			}
		}
	}
	for _, s := range d.Services {
		for _, n := range s.InputTypes {
			if n == "time.Time" || n == "sql.NullTime" {
				return true
			}
		}
		for _, n := range s.Output {
			if n == "time.Time" || n == "sql.NullTime" {
				return true
			}
		}
	}
	return false
}

func (d *Definition) importWrappers() bool {
	for _, m := range d.Messages {
		for _, typ := range m.AttrTypes {
			if strings.HasPrefix(typ, "sql.Null") && typ != "sql.NullTime" {
				return true
			}
		}
	}
	for _, s := range d.Services {
		for _, n := range s.InputTypes {
			if strings.HasPrefix(n, "sql.Null") && n != "sql.NullTime" {
				return true
			}
		}
		for _, n := range s.Output {
			if strings.HasPrefix(n, "sql.Null") && n != "sql.NullTime" {
				return true
			}
		}
	}
	return false
}

func ParseDefinition(src, engine, module string) (*Definition, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, src, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	if total := len(pkgs); total != 1 {
		return nil, fmt.Errorf("too many packages: %d", total)
	}

	for pkgName, pkg := range pkgs {
		def := Definition{
			Database: engine,
			Package:  pkgName,
			GoModule: module,
			SrcPath:  src,
			Messages: make(map[string]*Message),
		}

		for _, file := range pkg.Files {
			if file.Scope != nil {
				for name, obj := range file.Scope.Objects {
					if name == "Queries" {
						continue
					}
					if typ, ok := obj.Decl.(*ast.TypeSpec); ok {
						if structType, ok := typ.Type.(*ast.StructType); ok {
							def.Messages[name] = createMessage(name, structType)
						}
					}
				}
			}
			for _, n := range file.Decls {
				if fun, ok := n.(*ast.FuncDecl); ok {
					visitFunc(fun, &def)
				}
			}
		}

		return &def, nil
	}
	return nil, nil
}
