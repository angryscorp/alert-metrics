package main

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var exitCodeAnalyzer = &analysis.Analyzer{
	Name:     "exitcode",
	Doc:      "reports direct os.Exit calls in main functions",
	Run:      runExitCodeAnalyzer,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func runExitCodeAnalyzer(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
		(*ast.FuncDecl)(nil),
	}

	var inMainFunc bool

	nodeInspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeInspector.Preorder(nodeFilter, func(node ast.Node) {
		pos := pass.Fset.Position(node.Pos())
		if strings.Contains(pos.Filename, "go-build") {
			return
		}

		switch n := node.(type) {
		case *ast.FuncDecl:
			inMainFunc = n.Name.Name == "main"
		case *ast.CallExpr:
			if !inMainFunc {
				return
			}

			if sel, ok := n.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok {
					if ident.Name == "os" && sel.Sel.Name == "Exit" {
						pass.Reportf(n.Pos(), "direct os.Exit call in main function")
					}
				}
			}
		}
	})

	return nil, nil
}
