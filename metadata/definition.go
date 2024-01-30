package metadata

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/emicklei/proto"
	"github.com/walterwanderley/sqlc-grpc/converter"
)

type Definition struct {
	Args               string
	GoModule           string
	Packages           []*Package
	MigrationPath      string
	MigrationLib       string
	LiteFS             bool
	Litestream         bool
	DistributedTracing bool
	Metric             bool
}

func (d *Definition) Validate() error {
	switch d.MigrationLib {
	case "goose", "migrate":
	default:
		return fmt.Errorf("invalid migration library %q , use goose or go-migrate", d.MigrationLib)
	}
	var engine, sqlPackage string
	countServices := 0
	for _, pkg := range d.Packages {
		if engine == "" {
			engine = pkg.Engine
		} else if engine != pkg.Engine {
			return fmt.Errorf("can't use different database engines, found %q and %q", engine, pkg.Engine)
		}

		if sqlPackage == "" {
			sqlPackage = pkg.SqlPackage
		} else if sqlPackage != pkg.SqlPackage {
			return fmt.Errorf("can't use different sql packages, found %q and %q", sqlPackage, pkg.SqlPackage)
		}

		countServices += len(pkg.Services)
	}

	if countServices == 0 {
		return fmt.Errorf("no services found")
	}

	return nil

}

func (d *Definition) Database() string {
	for _, p := range d.Packages {
		if p.Engine != "" {
			return p.Engine
		}
	}
	return ""
}

func (d *Definition) DatabaseDriver() string {
	switch d.Database() {
	case "sqlite":
		if !d.LiteFS && !d.Litestream {
			return "sqlite"
		}
		return "sqlite3"
	case "postgresql":
		return "pgx"
	case "mysql":
		return "mysql"
	}

	return "unknown_driver"
}

func (d *Definition) DatabaseImport() string {
	switch d.Database() {
	case "sqlite":
		if !d.LiteFS && !d.Litestream {
			return "modernc.org/sqlite"
		}
		return "github.com/mattn/go-sqlite3"
	case "postgresql":
		if d.SqlPackage() == "pgx/v5" {
			return "github.com/jackc/pgx/v5/pgxpool"
		}
		return "github.com/jackc/pgx/v5/stdlib"
	case "mysql":
		return "github.com/go-sql-driver/mysql"
	}

	return "unknown_database"
}

func (d *Definition) SqlPackage() string {
	for _, p := range d.Packages {
		if p.SqlPackage != "" {
			return p.SqlPackage
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
	SqlPackage                 string
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
	CustomProtoRPCs            []string
	CustomProtoMessages        []string
	HasExecResult              bool
}

func (p *Package) ProtoImports() []string {
	r := make([]string, 0)
	if p.importTimestamp() {
		r = append(r, `import "google/protobuf/timestamp.proto";`)
	}
	if p.importWrappers() {
		r = append(r, `import "google/protobuf/wrappers.proto";`)
	}
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
		log.Println(err.Error())
		return
	}
	defer f.Close()
	parser := proto.NewParser(f)
	def, err := parser.Parse()
	if err != nil {
		log.Println(err.Error())
		return
	}

	proto.Walk(def, proto.WithImport(func(i *proto.Import) {
		if strings.Contains(i.Filename, "google/api/annotations.proto") ||
			strings.Contains(i.Filename, "protoc-gen-openapiv2/options/annotations.proto") {
			return
		}
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
		if s.Name != converter.UpperFirstCharacter(p.Package)+"Service" {
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
		options := make([]string, 0)

		for _, e := range rpc.Elements {
			opt, ok := e.(*proto.Option)
			if !ok {
				continue
			}
			options = append(options, fmt.Sprintf("option %s = {", opt.Name))
			options = append(options, printProtoLiteral(opt.Constant.OrderedMap, 1)...)
			options = append(options, "};")
		}
		var comments []string
		if rpc.Comment != nil {
			comments = clearLines(rpc.Comment.Lines)
		}
		var exists bool
		for _, s := range p.Services {
			if s.Name == rpc.Name {
				s.CustomProtoOptions = options
				s.CustomProtoComments = comments
				exists = true
				break
			}
		}
		if exists {
			return
		}

		for _, c := range comments {
			p.CustomProtoRPCs = append(p.CustomProtoRPCs, fmt.Sprintf("// %s", c))
		}
		var rpcInterface strings.Builder
		rpcInterface.WriteString("rpc ")
		rpcInterface.WriteString(rpc.Name)
		rpcInterface.WriteString("(")
		if rpc.StreamsRequest {
			rpcInterface.WriteString("streams ")
		}
		rpcInterface.WriteString(rpc.RequestType)
		rpcInterface.WriteString(") returns (")
		if rpc.StreamsReturns {
			rpcInterface.WriteString("streams ")
		}
		rpcInterface.WriteString(rpc.ReturnsType)
		rpcInterface.WriteString(")")
		if len(options) == 0 {
			rpcInterface.WriteString(";")
			p.CustomProtoRPCs = append(p.CustomProtoRPCs, rpcInterface.String())
		} else {
			rpcInterface.WriteString("{")
			p.CustomProtoRPCs = append(p.CustomProtoRPCs, rpcInterface.String())
			p.CustomProtoRPCs = append(p.CustomProtoRPCs, options...)
			p.CustomProtoRPCs = append(p.CustomProtoRPCs, "};")
		}
	}))

	proto.Walk(def, proto.WithMessage(func(protoMessage *proto.Message) {
		msg, ok := p.Messages[protoMessage.Name]
		if !ok {
			if strings.HasSuffix(protoMessage.Name, "Request") {
				msg, ok = p.Messages[protoMessage.Name[0:len(protoMessage.Name)-7]+"Params"]
				if !ok {
					p.addUserDefinedMessage(protoMessage)
					return
				}
			} else {
				p.addUserDefinedMessage(protoMessage)
				return
			}
		}
		msg.loadOptions(protoMessage)
	}))
}

func (p *Package) addUserDefinedMessage(protoMessage *proto.Message) {
	if protoMessage.Comment != nil {
		for _, c := range clearLines(protoMessage.Comment.Lines) {
			p.CustomProtoMessages = append(p.CustomProtoMessages, fmt.Sprintf("// %s", c))
		}
	}
	var options, fields []string
	for _, e := range protoMessage.Elements {
		if opt, ok := e.(*proto.Option); ok {
			options = append(options, fmt.Sprintf("    option %s = {", opt.Name))
			options = append(options, printProtoLiteral(opt.Constant.OrderedMap, 2)...)
			options = append(options, "    };")
			continue
		}

		if f, ok := e.(*proto.NormalField); ok {
			if f.Comment != nil {
				for _, c := range clearLines(f.Comment.Lines) {
					fields = append(fields, fmt.Sprintf("// %s", c))
				}
			}
			var fieldSpec strings.Builder
			if f.Repeated {
				fieldSpec.WriteString("repeated ")
			}
			fieldSpec.WriteString(f.Type)
			fieldSpec.WriteString(" ")
			fieldSpec.WriteString(f.Name)
			fieldSpec.WriteString(" = ")
			fieldSpec.WriteString(fmt.Sprintf("%d", f.Sequence))

			var fieldOptions []string
			var hasComplexOption bool
			for i, opt := range f.Options {
				var prefix string
				if i > 0 {
					prefix = "        "
				}
				var suffix string
				if i+1 < len(f.Options) {
					suffix = ", "
				}
				if opt.Constant.Source != "" {
					if hasComplexOption {
						fieldOptions = append(fieldOptions, fmt.Sprintf("%s%s = %s%s", prefix, opt.Name, opt.Constant.Source, suffix))
					} else {
						fieldOptions = append(fieldOptions, fmt.Sprintf("%s = %s%s", opt.Name, opt.Constant.Source, suffix))
					}
					continue
				}
				hasComplexOption = true
				fieldOptions = append(fieldOptions, fmt.Sprintf("%s%s = {\n", prefix, opt.Name))
				fieldOptions = append(fieldOptions, printProtoLiteral(opt.Constant.OrderedMap, 3)...)
				fieldOptions = append(fieldOptions, fmt.Sprintf("        }%s", suffix))
			}
			if len(fieldOptions) == 0 {
				fieldSpec.WriteString(";")
				fields = append(fields, fieldSpec.String())
				continue
			}
			fieldSpec.WriteString(" [")
			fields = append(fields, fieldSpec.String())
			fields = append(fields, fieldOptions...)
			fields = append(fields, "];")

		}
	}
	p.CustomProtoMessages = append(p.CustomProtoMessages, fmt.Sprintf("message %s {", protoMessage.Name))
	p.CustomProtoMessages = append(p.CustomProtoMessages, options...)
	p.CustomProtoMessages = append(p.CustomProtoMessages, fields...)
	p.CustomProtoMessages = append(p.CustomProtoMessages, "}")

}

func (p *Package) importTimestamp() bool {
	for _, m := range p.Messages {
		for _, f := range m.Fields {
			if f.Type == "time.Time" || f.Type == "sql.NullTime" || strings.HasPrefix(f.Type, "pgtype.Time") || strings.HasPrefix(f.Type, "pgtype.Date") {
				return true
			}
		}
	}
	for _, s := range p.Services {
		for _, n := range s.InputTypes {
			if n == "time.Time" || n == "sql.NullTime" || strings.HasPrefix(n, "pgtype.Time") || strings.HasPrefix(n, "pgtype.Date") {
				return true
			}
		}

		if s.Output == "time.Time" || s.Output == "sql.NullTime" || strings.HasPrefix(s.Output, "pgtype.Time") || strings.HasPrefix(s.Output, "pgtype.Date") {
			return true
		}

	}
	return false
}

func (p *Package) importWrappers() bool {
	for _, m := range p.Messages {
		for _, f := range m.Fields {
			if (strings.HasPrefix(f.Type, "sql.Null") || strings.HasPrefix(f.Type, "pgtype.")) && f.Type != "sql.NullTime" {
				return true
			}
		}
	}
	for _, s := range p.Services {
		for _, n := range s.InputTypes {
			if (strings.HasPrefix(n, "sql.Null") || strings.HasPrefix(n, "pgtype.")) && n != "sql.NullTime" {
				return true
			}
		}

		if (strings.HasPrefix(s.Output, "sql.Null") || strings.HasPrefix(s.Output, "pgtype.")) && s.Output != "sql.NullTime" {
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

		}

		for _, m := range p.Messages {
			m.adjustType(p.Messages)
		}

		for _, file := range pkg.Files {
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

		sort.SliceStable(p.Services, func(i, j int) bool {
			return strings.Compare(p.Services[i].Name, p.Services[j].Name) < 0
		})

		outAdapters := make(map[string]struct{})

		for _, s := range p.Services {
			if s.HasCustomOutput() || s.HasArrayOutput() {
				outAdapters[converter.CanonicalName(s.Output)] = struct{}{}
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
			constants[name] = v.Value
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
