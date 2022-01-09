package metadata

import (
	"go/ast"
	"strings"
)

func visitFunc(fun *ast.FuncDecl, def *Package, constants map[string]string) {
	if !isMethodValid(fun) {
		return
	}

	inputNames := make([]string, 0)
	inputTypes := make([]string, 0)

	// context is the first parameter
	for i := 1; i < len(fun.Type.Params.List); i++ {
		p := fun.Type.Params.List[i]
		for _, n := range p.Names {
			typ, err := exprToStr(p.Type)
			if err != nil {
				return
			}
			if typ == "DBTX" {
				continue
			}
			inputTypes = append(inputTypes, typ)
			inputNames = append(inputNames, n.Name)
		}
	}

	var output string
	// two is the maximum results for a valid method, error is the last result
	if len(fun.Type.Results.List) > 1 {
		p := fun.Type.Results.List[0]
		var err error
		output, err = exprToStr(p.Type)
		if err != nil {
			return
		}

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

	t, ok := typ.X.(*ast.Ident)
	if !ok {
		return false
	}

	if t.Name != "Queries" {
		return false
	}

	firstParam, err := exprToStr(fun.Type.Params.List[0].Type)
	if err != nil {
		return false
	}

	if firstParam != "context.Context" {
		return false
	}

	if len(fun.Type.Results.List) > 2 {
		return false
	}

	lastResult, err := exprToStr(fun.Type.Results.List[len(fun.Type.Results.List)-1].Type)
	if err != nil {
		return false
	}

	if lastResult != "error" {
		return false
	}

	return true
}

func canonicalName(typ string) string {
	name := strings.TrimPrefix(typ, "[]")
	name = strings.TrimPrefix(name, "*")
	return name
}
