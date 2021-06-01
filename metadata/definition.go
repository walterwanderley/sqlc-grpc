package metadata

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type Definition struct {
	GoModule string
	Packages []*Package
}

func (d *Definition) Database() string {
	for _, p := range d.Packages {
		if p.Engine != "" {
			return p.Engine
		}
	}
	return ""
}

type Package struct {
	Engine   string
	Package  string
	GoModule string
	SrcPath  string
	Services []*Service
	Messages map[string]*Message
}

func (p *Package) ProtoImports() []string {
	r := make([]string, 0)
	if p.importEmpty() {
		r = append(r, `import "google/protobuf/empty.proto";`)
	}
	if p.importTimestamp() {
		r = append(r, `import "google/protobuf/timestamp.proto";`)
	}
	if p.importWrappers() {
		r = append(r, `import "google/protobuf/wrappers.proto";`)
	}
	return r
}

func (p *Package) importEmpty() bool {
	for _, s := range p.Services {
		if s.EmptyInput() || s.EmptyOutput() {
			return true
		}
	}
	return false
}

func (p *Package) importTimestamp() bool {
	for _, m := range p.Messages {
		for _, typ := range m.AttrTypes {
			if typ == "time.Time" || typ == "sql.NullTime" {
				return true
			}
		}
	}
	for _, s := range p.Services {
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

func (p *Package) importWrappers() bool {
	for _, m := range p.Messages {
		for _, typ := range m.AttrTypes {
			if strings.HasPrefix(typ, "sql.Null") && typ != "sql.NullTime" {
				return true
			}
		}
	}
	for _, s := range p.Services {
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

func ParsePackage(src string) (*Package, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, src, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	if total := len(pkgs); total != 1 {
		return nil, fmt.Errorf("too many packages: %d", total)
	}

	for pkgName, pkg := range pkgs {
		p := Package{
			Package:  pkgName,
			SrcPath:  src,
			Messages: make(map[string]*Message),
		}

		for _, file := range pkg.Files {
			if file.Scope != nil {
				for name, obj := range file.Scope.Objects {
					if name == "Queries" || name == "Service" {
						continue
					}
					if typ, ok := obj.Decl.(*ast.TypeSpec); ok {
						if structType, ok := typ.Type.(*ast.StructType); ok {
							p.Messages[name] = createMessage(name, structType)
						}
					}
				}
			}
			for _, n := range file.Decls {
				if fun, ok := n.(*ast.FuncDecl); ok {
					visitFunc(fun, &p)
				}
			}
		}

		return &p, nil
	}
	return nil, nil
}