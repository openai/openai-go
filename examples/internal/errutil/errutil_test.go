package errutil_test

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/openai/openai-go/examples/internal/errutil"
	"github.com/openai/openai-go/v3"
)

func TestMessage(t *testing.T) {
	apiErr := &openai.Error{StatusCode: 429, Message: "synthetic-secret-message", Code: "synthetic-secret-code"}
	for _, tc := range []struct {
		name string
		err  error
		want string
	}{
		{"API", apiErr, "OpenAI API request failed (HTTP 429 Too Many Requests)"},
		{"wrapped API", fmt.Errorf("synthetic-secret-wrapper: %w", apiErr), "OpenAI API request failed (HTTP 429 Too Many Requests)"},
		{"transport", &url.Error{Op: "Post", URL: "https://synthetic-secret.invalid/?token=synthetic-secret", Err: errors.New("synthetic-secret")}, "Example operation failed"},
		{"provider", errors.New("synthetic-secret-provider-error"), "Example operation failed"},
		{"nil API", (*openai.Error)(nil), "Example operation failed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := errutil.Message(tc.err); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
	// Unknown HTTP statuses must not cause arbitrary response status text to be used.
	if got := errutil.Message(&openai.Error{StatusCode: 499, Response: &http.Response{Status: "synthetic-secret"}}); strings.Contains(got, "synthetic-secret") {
		t.Fatal(got)
	}
}

// Keep every runnable example's terminal error sinks on the tested safe policy.
// This includes cloud and mutual-TLS examples that cannot run without credentials.
func TestExampleErrorSinks(t *testing.T) {
	paths, err := filepath.Glob("../../*/main.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) == 0 {
		t.Fatal("no runnable examples found")
	}
	for _, path := range paths {
		t.Run(filepath.Base(filepath.Dir(path)), func(t *testing.T) {
			fs := token.NewFileSet()
			file, err := parser.ParseFile(fs, path, nil, 0)
			if err != nil {
				t.Fatal(err)
			}
			sinks := 0
			ast.Inspect(file, func(node ast.Node) bool {
				call, ok := node.(*ast.CallExpr)
				if !ok {
					return true
				}
				var args []ast.Expr
				if id, ok := call.Fun.(*ast.Ident); ok && id.Name == "panic" {
					args = call.Args
				}
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					if pkg, ok := sel.X.(*ast.Ident); ok {
						if pkg.Name == "log" && sel.Sel.Name == "Fatal" {
							args = call.Args
						}
						if pkg.Name == "slog" && sel.Sel.Name == "Error" {
							args = call.Args
						}
					}
				}
				if len(args) == 0 {
					return true
				}
				sinks++
				for _, arg := range args {
					if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
						if _, err := strconv.Unquote(lit.Value); err == nil {
							continue
						}
					}
					if safeCall, ok := arg.(*ast.CallExpr); ok {
						if sel, ok := safeCall.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "Message" {
							if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "errutil" {
								continue
							}
						}
					}
					t.Errorf("%s: error sink must use a constant or errutil.Message", fs.Position(arg.Pos()))
				}
				return true
			})
			if sinks == 0 {
				t.Fatal("example has no checked terminal error sink")
			}
		})
	}
}

func TestStatus(t *testing.T) {
	for _, status := range []string{"processed", "running", "failed", "completed"} {
		if got := errutil.Status(status); got != status {
			t.Fatalf("known status %q changed to %q", status, got)
		}
	}
	if got := errutil.Status("synthetic-secret-status"); got != "unknown" {
		t.Fatalf("unknown status exposed: %q", got)
	}
}
