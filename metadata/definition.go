package metadata

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/emicklei/proto"
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
	Engine                     string
	Package                    string
	GoModule                   string
	SchemaPath                 string
	SrcPath                    string
	Services                   []*Service
	Messages                   map[string]*Message
	OutputAdapters             []*Message
	EmitInterface              bool
	EmitParamsPointers         bool
	EmitResultPointers         bool
	EmitDbArgument             bool
	CustomProtoOptions         []string
	CustomProtoImports         []string
	CustomServiceProtoComments []string
	CustomServiceProtoOptions  []string
}

func (p *Package) ProtoImports() []string {
	r := make([]string, 0)
	r = append(r, `import "google/api/annotations.proto";`)
	if p.importTimestamp() {
		r = append(r, `import "google/protobuf/timestamp.proto";`)
	}
	if p.importWrappers() {
		r = append(r, `import "google/protobuf/wrappers.proto";`)
	}
	r = append(r, `import "protoc-gen-openapiv2/options/annotations.proto";`)
	imports := strings.Join(r, " ")
	for _, i := range p.CustomProtoImports {
		if !strings.Contains(imports, i) {
			r = append(r, fmt.Sprintf("import \"%s\";", i))
		}
	}
	return r
}

func (p *Package) LoadOptions(protoFile string) {
	f, err := os.Open(protoFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()
	parser := proto.NewParser(f)
	def, err := parser.Parse()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	proto.Walk(def, proto.WithImport(func(i *proto.Import) {
		p.CustomProtoImports = append(p.CustomProtoImports, i.Filename)
	}))

	proto.Walk(def, proto.WithOption(func(opt *proto.Option) {
		if _, ok := opt.Parent.(*proto.Proto); ok {
			if opt.Constant.Source != "" {
				p.CustomProtoOptions = append(p.CustomProtoOptions, fmt.Sprintf("option %s = \"%s\";", opt.Name, opt.Constant.Source))
			} else {
				p.CustomProtoOptions = append(p.CustomProtoOptions, fmt.Sprintf("option %s = {", opt.Name))
				p.CustomProtoOptions = append(p.CustomProtoOptions, printProtoLiteral(opt.Constant.OrderedMap, 1)...)
				p.CustomProtoOptions = append(p.CustomProtoOptions, "};")
			}
		}
	}))

	proto.Walk(def, proto.WithService(func(s *proto.Service) {
		if s.Name != UpperFirstCharacter(p.Package)+"Service" {
			return
		}
		if s.Comment != nil {
			p.CustomServiceProtoComments = clearLines(s.Comment.Lines)
		}
		for _, e := range s.Elements {
			opt, ok := e.(*proto.Option)
			if !ok {
				continue
			}
			p.CustomServiceProtoOptions = append(p.CustomServiceProtoOptions, fmt.Sprintf("option %s = {", opt.Name))
			p.CustomServiceProtoOptions = append(p.CustomServiceProtoOptions, printProtoLiteral(opt.Constant.OrderedMap, 1)...)
			p.CustomServiceProtoOptions = append(p.CustomServiceProtoOptions, "};")
		}

	}))

	proto.Walk(def, proto.WithRPC(func(rpc *proto.RPC) {
		res := make([]string, 0)

		for _, e := range rpc.Elements {
			opt, ok := e.(*proto.Option)
			if !ok {
				continue
			}
			res = append(res, fmt.Sprintf("option %s = {", opt.Name))
			res = append(res, printProtoLiteral(opt.Constant.OrderedMap, 1)...)
			res = append(res, "};")
		}

		for _, s := range p.Services {
			if s.Name == rpc.Name {
				s.CustomProtoOptions = res
				if rpc.Comment != nil {
					s.CustomProtoComments = clearLines(rpc.Comment.Lines)
				}
				break
			}
		}
	}))

	proto.Walk(def, proto.WithMessage(func(protoMessage *proto.Message) {
		msg, ok := p.Messages[protoMessage.Name]
		if !ok {
			if strings.HasSuffix(protoMessage.Name, "Request") {
				msg, ok = p.Messages[protoMessage.Name[0:len(protoMessage.Name)-7]+"Params"]
				if !ok {
					return
				}
			} else {
				return
			}
		}
		msg.loadOptions(protoMessage)
	}))
}

func (p *Package) importTimestamp() bool {
	for _, m := range p.Messages {
		for _, f := range m.Fields {
			if f.Type == "time.Time" || f.Type == "sql.NullTime" {
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

		if s.Output == "time.Time" || s.Output == "sql.NullTime" {
			return true
		}

	}
	return false
}

func (p *Package) importWrappers() bool {
	for _, m := range p.Messages {
		for _, f := range m.Fields {
			if strings.HasPrefix(f.Type, "sql.Null") && f.Type != "sql.NullTime" {
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

		if strings.HasPrefix(s.Output, "sql.Null") && s.Output != "sql.NullTime" {
			return true
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

		outAdapters := make(map[string]struct{})

		for _, s := range p.Services {
			if s.HasCustomOutput() || s.HasArrayOutput() {
				outAdapters[canonicalName(s.Output)] = struct{}{}
			}
		}

		p.OutputAdapters = make([]*Message, len(outAdapters))
		i := 0
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

func printProtoLiteral(literal proto.LiteralMap, deep int) []string {
	res := make([]string, 0)
	layout := fmt.Sprintf("%%-%ds", deep*4)
	prefix := fmt.Sprintf(layout, "")
	for _, item := range literal {
		if item.IsString {
			res = append(res, fmt.Sprintf("%s%s: \"%s\"", prefix, item.Name, item.Source))
		} else {
			if len(item.Array) > 0 {
				items := make([]string, 0)
				for _, i := range item.Array {
					items = append(items, fmt.Sprintf(`"%s"`, i.Source))
				}
				res = append(res, fmt.Sprintf("%s%s: [%s]", prefix, item.Name, strings.Join(items, ", ")))
			} else {
				res = append(res, fmt.Sprintf("%s%s: {", prefix, item.Name))
				res = append(res, printProtoLiteral(item.OrderedMap, deep+1)...)
				res = append(res, fmt.Sprintf("%s};", prefix))
			}
		}
	}
	return res
}

func clearLines(lines []string) []string {
	res := make([]string, 0)
	for _, l := range lines {
		line := strings.TrimSpace(l)
		if len(line) > 0 {
			res = append(res, line)
		}
	}
	return res
}
