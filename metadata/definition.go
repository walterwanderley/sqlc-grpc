package metadata

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"sort"
	"strings"
)

type Definition struct {
	Args     string
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

type PackageOpts struct {
	Path               string
	EmitInterface      bool
	EmitParamsPointers bool
	EmitResultPointers bool
	EmitDbArgument     bool
}

type Package struct {
	Engine             string
	Package            string
	GoModule           string
	SchemaPath         string
	SrcPath            string
	Services           []*Service
	Messages           map[string]*Message
	InputAdapters      []*Message
	OutputAdapters     []*Message
	EmitInterface      bool
	EmitParamsPointers bool
	EmitResultPointers bool
	EmitDbArgument     bool
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

func ParsePackage(opts PackageOpts, queriesToIgnore []*regexp.Regexp) (*Package, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, opts.Path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	if total := len(pkgs); total != 1 {
		return nil, fmt.Errorf("too many packages: %d", total)
	}

	for pkgName, pkg := range pkgs {
		p := Package{
			Package:            pkgName,
			SrcPath:            opts.Path,
			Messages:           make(map[string]*Message),
			EmitInterface:      opts.EmitInterface,
			EmitParamsPointers: opts.EmitParamsPointers,
			EmitResultPointers: opts.EmitResultPointers,
			EmitDbArgument:     opts.EmitDbArgument,
		}

		constants := make(map[string]string)
		for _, file := range pkg.Files {
			if file.Scope != nil {
				for name, obj := range file.Scope.Objects {
					if name == "Queries" || name == "Service" {
						continue
					}
					addConstant(constants, name, obj)
					if typ, ok := obj.Decl.(*ast.TypeSpec); ok {
						switch t := typ.Type.(type) {
						case *ast.Ident:
							msg, err := createAliasMessage(name, t)
							if err != nil {
								return nil, err
							}
							p.Messages[name] = msg
						case *ast.StructType:
							msg, err := createStructMessage(name, t)
							if err != nil {
								return nil, err
							}
							p.Messages[name] = msg
						case *ast.ArrayType:
							msg, err := createArrayMessage(name, t)
							if err != nil {
								return nil, err
							}
							p.Messages[name] = msg
						}
					}

				}
			}
			for _, n := range file.Decls {
				if fun, ok := n.(*ast.FuncDecl); ok {
					var ignore bool
					for _, re := range queriesToIgnore {
						if re.MatchString(fun.Name.String()) {
							ignore = true
							break
						}
					}
					if !ignore {
						visitFunc(fun, &p, constants)
					}
				}
			}
		}

		for _, m := range p.Messages {
			m.adjustType(p.Messages)
		}

		sort.SliceStable(p.Services, func(i, j int) bool {
			return strings.Compare(p.Services[i].Name, p.Services[j].Name) < 0
		})

		inAdapters := make(map[string]struct{})
		outAdapters := make(map[string]struct{})

		for _, s := range p.Services {
			if s.HasCustomParams() {
				inAdapters[canonicalName(s.InputTypes[0])] = struct{}{}
			}
			if s.HasCustomOutput() {
				for _, n := range s.Output {
					outAdapters[canonicalName(n)] = struct{}{}
				}
			}
			if s.HasArrayOutput() {
				outAdapters[canonicalName(s.Output[0])] = struct{}{}
			}
		}

		p.InputAdapters = make([]*Message, len(inAdapters))
		i := 0
		for k := range inAdapters {
			p.InputAdapters[i] = p.Messages[k]
			i++
		}

		sort.SliceStable(p.InputAdapters, func(i, j int) bool {
			return strings.Compare(p.InputAdapters[i].Name, p.InputAdapters[j].Name) < 0
		})

		p.OutputAdapters = make([]*Message, len(outAdapters))
		i = 0
		for k := range outAdapters {
			p.OutputAdapters[i] = p.Messages[k]
			i++
		}

		sort.SliceStable(p.OutputAdapters, func(i, j int) bool {
			return strings.Compare(p.OutputAdapters[i].Name, p.OutputAdapters[j].Name) < 0
		})

		return &p, nil
	}
	return nil, nil
}

func addConstant(constants map[string]string, name string, obj *ast.Object) {
	if obj.Kind != ast.Con {
		return
	}
	if vs, ok := obj.Decl.(*ast.ValueSpec); ok {
		if v, ok := vs.Values[0].(*ast.BasicLit); ok {
			constants[UpperFirstCharacter(name)] = v.Value
		}
	}

}
