package config

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type GoConfig struct {
	editor EditorAPI
}

type EditorAPI struct {
	BindKey       func(key, command string)
	LoadPlugin    func(name string)
	SetOption     func(key string, value any)
	RegisterHook  func(event string, handler func())
	InstallPlugin func(name, repository, version string)
}

func LoadGoConfig(filepath string, api EditorAPI) error {
	goConfig := &GoConfig{editor: api}
	
	src, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filepath, src, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse Go config: %v", err)
	}
	
	return goConfig.executeAST(f)
}

func (gc *GoConfig) executeAST(f *ast.File) error {
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if fn.Name.Name == "init" {
				if err := gc.executeBlock(fn.Body); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (gc *GoConfig) executeBlock(block *ast.BlockStmt) error {
	for _, stmt := range block.List {
		if err := gc.executeStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (gc *GoConfig) executeStmt(stmt ast.Stmt) error {
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		return gc.executeExpr(s.X)
	case *ast.IfStmt:
		return gc.executeIfStmt(s)
	case *ast.ForStmt:
		return gc.executeForStmt(s)
	case *ast.AssignStmt:
		return gc.executeAssignStmt(s)
	}
	return nil
}

func (gc *GoConfig) executeExpr(expr ast.Expr) error {
	switch e := expr.(type) {
	case *ast.CallExpr:
		return gc.executeCallExpr(e)
	}
	return nil
}

func (gc *GoConfig) executeCallExpr(call *ast.CallExpr) error {
	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:
		if ident, ok := fun.X.(*ast.Ident); ok && ident.Name == "editor" {
			return gc.executeEditorCall(fun.Sel.Name, call.Args)
		}
	case *ast.Ident:
		switch fun.Name {
		case "BindKey":
			return gc.executeBindKey(call.Args)
		case "LoadPlugin":
			return gc.executeLoadPlugin(call.Args)
		case "SetOption":
			return gc.executeSetOption(call.Args)
		case "InstallPlugin":
			return gc.executeInstallPlugin(call.Args)
		}
	}
	return nil
}

func (gc *GoConfig) executeEditorCall(method string, args []ast.Expr) error {
	switch method {
	case "BindKey":
		return gc.executeBindKey(args)
	case "LoadPlugin":
		return gc.executeLoadPlugin(args)
	case "SetOption":
		return gc.executeSetOption(args)
	case "RegisterHook":
		return gc.executeRegisterHook(args)
	case "InstallPlugin":
		return gc.executeInstallPlugin(args)
	}
	return nil
}

func (gc *GoConfig) executeBindKey(args []ast.Expr) error {
	if len(args) < 2 {
		return fmt.Errorf("BindKey requires 2 arguments")
	}
	
	key, err := gc.evalStringLiteral(args[0])
	if err != nil {
		return err
	}
	
	command, err := gc.evalStringLiteral(args[1])
	if err != nil {
		return err
	}
	
	gc.editor.BindKey(key, command)
	return nil
}

func (gc *GoConfig) executeLoadPlugin(args []ast.Expr) error {
	if len(args) < 1 {
		return fmt.Errorf("LoadPlugin requires 1 argument")
	}
	
	name, err := gc.evalStringLiteral(args[0])
	if err != nil {
		return err
	}
	
	gc.editor.LoadPlugin(name)
	return nil
}

func (gc *GoConfig) executeSetOption(args []ast.Expr) error {
	if len(args) < 2 {
		return fmt.Errorf("SetOption requires 2 arguments")
	}
	
	key, err := gc.evalStringLiteral(args[0])
	if err != nil {
		return err
	}
	
	value, err := gc.evalExpression(args[1])
	if err != nil {
		return err
	}
	
	gc.editor.SetOption(key, value)
	return nil
}

func (gc *GoConfig) executeRegisterHook(args []ast.Expr) error {
	if len(args) < 2 {
		return fmt.Errorf("RegisterHook requires 2 arguments")
	}
	
	event, err := gc.evalStringLiteral(args[0])
	if err != nil {
		return err
	}
	
	gc.editor.RegisterHook(event, func() {
	})
	return nil
}

func (gc *GoConfig) executeInstallPlugin(args []ast.Expr) error {
	if len(args) < 3 {
		return fmt.Errorf("InstallPlugin requires 3 arguments")
	}
	
	name, err := gc.evalStringLiteral(args[0])
	if err != nil {
		return err
	}
	
	repository, err := gc.evalStringLiteral(args[1])
	if err != nil {
		return err
	}
	
	version, err := gc.evalStringLiteral(args[2])
	if err != nil {
		return err
	}
	
	if gc.editor.InstallPlugin != nil {
		gc.editor.InstallPlugin(name, repository, version)
	}
	return nil
}

func (gc *GoConfig) executeIfStmt(stmt *ast.IfStmt) error {
	return nil
}

func (gc *GoConfig) executeForStmt(stmt *ast.ForStmt) error {
	return nil
}

func (gc *GoConfig) executeAssignStmt(stmt *ast.AssignStmt) error {
	return nil
}

func (gc *GoConfig) evalStringLiteral(expr ast.Expr) (string, error) {
	if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		return strings.Trim(lit.Value, "\""), nil
	}
	return "", fmt.Errorf("expected string literal")
}

func (gc *GoConfig) evalExpression(expr ast.Expr) (any, error) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		switch e.Kind {
		case token.STRING:
			return strings.Trim(e.Value, "\""), nil
		case token.INT:
			return strconv.Atoi(e.Value)
		case token.FLOAT:
			return strconv.ParseFloat(e.Value, 64)
		}
	case *ast.Ident:
		switch e.Name {
		case "true":
			return true, nil
		case "false":
			return false, nil
		}
	}
	return nil, fmt.Errorf("unsupported expression type: %s", reflect.TypeOf(expr))
}