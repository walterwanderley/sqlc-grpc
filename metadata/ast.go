package metadata

import "go/ast"

func visitFunc(fun *ast.FuncDecl, def *Package, constants map[string]string) {
	if !isMethodValid(fun) {
		return
	}

	inputNames := make([]string, 0)
	inputTypes := make([]string, 0)
	output := make([]string, 0)

	// context is the first parameter
	for i := 1; i < len(fun.Type.Params.List); i++ {
		p := fun.Type.Params.List[i]
		inputNames = append(inputNames, p.Names[0].Name)
		inputTypes = append(inputTypes, exprToStr(p.Type))
	}

	// error is the last result
	for i := 0; i < len(fun.Type.Results.List)-1; i++ {
		p := fun.Type.Results.List[0]
		output = append(output, exprToStr(p.Type))
	}
	def.Services = append(def.Services, &Service{
		Name:       fun.Name.String(),
		InputNames: inputNames,
		InputTypes: inputTypes,
		Output:     output,
		Sql:        constants[fun.Name.String()],
		Messages:   def.Messages,
	})
}

func createMessage(name string, s *ast.StructType) *Message {
	names := make([]string, 0)
	types := make([]string, 0)
	for _, f := range s.Fields.List {
		types = append(types, exprToStr(f.Type))
		var name string
		if len(f.Names) > 0 {
			name = f.Names[0].Name
		}
		names = append(names, name)
	}
	return &Message{
		Name:      name,
		AttrNames: names,
		AttrTypes: types,
	}
}

func isMethodValid(fun *ast.FuncDecl) bool {
	if fun.Name == nil {
		return false
	}

	if !fun.Name.IsExported() {
		return false
	}

	if fun.Recv == nil || len(fun.Recv.List) != 1 {
		return false
	}

	typ, ok := fun.Recv.List[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}

	if fun.Type.Params == nil || len(fun.Type.Params.List) == 0 ||
		fun.Type.Results == nil || len(fun.Type.Results.List) == 0 {
		return false
	}

	if exprToStr(fun.Type.Params.List[0].Type) != "context.Context" {
		return false
	}

	if exprToStr(fun.Type.Results.List[len(fun.Type.Results.List)-1].Type) != "error" {
		return false
	}

	t, ok := typ.X.(*ast.Ident)
	if !ok {
		return false
	}

	return t.Name == "Queries"
}
