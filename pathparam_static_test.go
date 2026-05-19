package openai_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

type requestPathKind int

const (
	requestPathLiteral requestPathKind = 1 << iota
	requestPathFormatPath
	requestPathUnknown
)

const requestconfigImportPath = "github.com/openai/openai-go/v3/internal/requestconfig"

func TestGeneratedPathsUseFormatPath(t *testing.T) {
	var violations []string
	fset := token.NewFileSet()

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", ".github", ".codex", "examples":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		source, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !isGeneratedServiceSource(path, source) {
			return nil
		}

		file, err := parser.ParseFile(fset, path, source, 0)
		if err != nil {
			return err
		}
		requestconfigNames := requestconfigImportNames(file)
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			reviewRequestPathFunction(fset, fn, requestconfigNames, &violations)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(violations) > 0 {
		t.Fatalf("generated request paths with parameters must use requestconfig.FormatPath:\n%s", strings.Join(violations, "\n"))
	}
}

func isGeneratedServiceSource(path string, source []byte) bool {
	if path == "client.go" {
		return false
	}
	return strings.HasPrefix(string(source), "// File generated from our OpenAPI spec by Stainless.")
}

func requestconfigImportNames(file *ast.File) map[string]bool {
	names := map[string]bool{}
	for _, importSpec := range file.Imports {
		importPath, err := strconv.Unquote(importSpec.Path.Value)
		if err != nil || importPath != requestconfigImportPath {
			continue
		}
		if importSpec.Name != nil && importSpec.Name.Name != "." && importSpec.Name.Name != "_" {
			names[importSpec.Name.Name] = true
			continue
		}
		names["requestconfig"] = true
	}
	return names
}

func reviewRequestPathFunction(fset *token.FileSet, fn *ast.FuncDecl, requestconfigNames map[string]bool, violations *[]string) {
	values := map[string]requestPathKind{}

	ast.Inspect(fn.Body, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.AssignStmt:
			for i, lhs := range n.Lhs {
				ident, ok := lhs.(*ast.Ident)
				if !ok || ident.Name == "_" || i >= len(n.Rhs) {
					continue
				}
				values[ident.Name] |= classifyRequestPathExpr(fset, n.Rhs[i], requestconfigNames, violations)
			}
		case *ast.ValueSpec:
			for i, name := range n.Names {
				if i >= len(n.Values) {
					continue
				}
				values[name.Name] |= classifyRequestPathExpr(fset, n.Values[i], requestconfigNames, violations)
			}
		case *ast.CallExpr:
			if !isRequestConfigCall(n, requestconfigNames, "ExecuteNewRequest") && !isRequestConfigCall(n, requestconfigNames, "NewRequestConfig") {
				return true
			}
			if len(n.Args) < 3 {
				*violations = append(*violations, fmt.Sprintf("%s: requestconfig call is missing request path argument", fset.Position(n.Pos())))
				return true
			}
			if !classifyRequestPathArg(fset, n.Args[2], values, requestconfigNames, violations).safe() {
				*violations = append(*violations, fmt.Sprintf("%s: dynamic request path must come from requestconfig.FormatPath", fset.Position(n.Args[2].Pos())))
			}
		}
		return true
	})
}

func classifyRequestPathArg(fset *token.FileSet, expr ast.Expr, values map[string]requestPathKind, requestconfigNames map[string]bool, violations *[]string) requestPathKind {
	switch e := expr.(type) {
	case *ast.Ident:
		return values[e.Name]
	default:
		return classifyRequestPathExpr(fset, expr, requestconfigNames, violations)
	}
}

func classifyRequestPathExpr(fset *token.FileSet, expr ast.Expr, requestconfigNames map[string]bool, violations *[]string) requestPathKind {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			if literalRequestPathHasPlaceholder(e) {
				return requestPathUnknown
			}
			return requestPathLiteral
		}
	case *ast.CallExpr:
		if isRequestConfigCall(e, requestconfigNames, "FormatPath") {
			validateFormatPathCall(fset, e, violations)
			return requestPathFormatPath
		}
	}
	return requestPathUnknown
}

func literalRequestPathHasPlaceholder(lit *ast.BasicLit) bool {
	value, err := strconv.Unquote(lit.Value)
	if err != nil {
		return true
	}
	return strings.Contains(value, "%s") || strings.Contains(value, "{") || strings.Contains(value, "}")
}

func (kind requestPathKind) safe() bool {
	return kind != 0 && kind&requestPathUnknown == 0
}

func validateFormatPathCall(fset *token.FileSet, call *ast.CallExpr, violations *[]string) {
	if len(call.Args) == 0 {
		*violations = append(*violations, fmt.Sprintf("%s: requestconfig.FormatPath is missing a format string", fset.Position(call.Pos())))
		return
	}

	formatLit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || formatLit.Kind != token.STRING {
		*violations = append(*violations, fmt.Sprintf("%s: requestconfig.FormatPath format must be a string literal", fset.Position(call.Args[0].Pos())))
		return
	}
	format, err := strconv.Unquote(formatLit.Value)
	if err != nil {
		*violations = append(*violations, fmt.Sprintf("%s: requestconfig.FormatPath format is not a valid string literal", fset.Position(formatLit.Pos())))
		return
	}

	placeholders, unsupported := countStringPlaceholders(format)
	if unsupported != "" {
		*violations = append(*violations, fmt.Sprintf("%s: requestconfig.FormatPath only supports %%s placeholders, found %%%s", fset.Position(formatLit.Pos()), unsupported))
	}
	if placeholders != len(call.Args)-1 {
		*violations = append(*violations, fmt.Sprintf("%s: requestconfig.FormatPath has %d %%s placeholders but %d path params", fset.Position(formatLit.Pos()), placeholders, len(call.Args)-1))
	}
}

func countStringPlaceholders(format string) (int, string) {
	count := 0
	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			continue
		}
		if i+1 >= len(format) {
			return count, "EOF"
		}
		i++
		switch format[i] {
		case '%':
			continue
		case 's':
			count++
		default:
			return count, string(format[i])
		}
	}
	return count, ""
}

func isRequestConfigCall(call *ast.CallExpr, requestconfigNames map[string]bool, name string) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != name {
		return false
	}
	pkg, ok := selector.X.(*ast.Ident)
	return ok && requestconfigNames[pkg.Name]
}
