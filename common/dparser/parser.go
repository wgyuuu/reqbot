package dparser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/davecgh/go-spew/spew"
)

func pkgFunc(src string) string {
	return `
    func Get() (list []int) {
        for i := 1; i < 10; i++ {
            list = append(list, i)
        }
        return
    }`
}

func Run(src string) (interface{}, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", pkgFunc(src), parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			break
		}

		if fn.Name.Name != "Get" {
			break
		}

		return execGet(fn)
	}

	return nil, errors.New("function name not Get")
}

func execGet(fn *ast.FuncDecl) (interface{}, error) {
	for _, stmt := range fn.Body.List {
		fmt.Println(spew.Sdump(stmt))
	}
	//
	return nil, nil
}
