package should

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestWrapperParityWithAssertPackage(t *testing.T) {
	t.Parallel()

	shouldFuncs := mustParseExportedFuncs(t, "should.go")
	assertFuncs := mustParseExportedFuncs(t, "assert/assertions.go")

	shared := 0
	for name, shouldDecl := range shouldFuncs {
		assertDecl, ok := assertFuncs[name]
		if !ok {
			continue
		}
		shared++

		shouldDoc := normalizedDoc(shouldDecl.Doc)
		assertDoc := normalizedDoc(assertDecl.Doc)
		if shouldDoc != assertDoc {
			t.Errorf("doc comment mismatch for %s\nshould.go:\n%s\nassert/assertions.go:\n%s", name, shouldDoc, assertDoc)
		}

		shouldSig := mustFormatFuncType(t, shouldDecl.Type)
		assertSig := mustFormatFuncType(t, assertDecl.Type)
		if shouldSig != assertSig {
			t.Errorf("signature mismatch for %s\nshould.go: %s\nassert/assertions.go: %s", name, shouldSig, assertSig)
		}
	}

	if shared == 0 {
		t.Fatal("no shared exported functions found between should.go and assert/assertions.go")
	}
}

func mustParseExportedFuncs(t *testing.T, path string) map[string]*ast.FuncDecl {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}

	funcs := make(map[string]*ast.FuncDecl)
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || !fn.Name.IsExported() {
			continue
		}
		funcs[fn.Name.Name] = fn
	}

	return funcs
}

func normalizedDoc(doc *ast.CommentGroup) string {
	if doc == nil {
		return ""
	}
	return strings.TrimSpace(doc.Text())
}

func mustFormatFuncType(t *testing.T, fnType *ast.FuncType) string {
	t.Helper()

	var buf bytes.Buffer
	if err := format.Node(&buf, token.NewFileSet(), fnType); err != nil {
		t.Fatalf("format function type: %v", err)
	}

	return strings.ReplaceAll(buf.String(), "assert.", "")
}
