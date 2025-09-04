package main

// Tool: infer_api_from_usage
// Scans packages/tui for usages of package "opencode" and gathers
// selector-based field accesses and callsites to populate api.md.

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"sort"

	"golang.org/x/tools/go/packages"
)

func main() {
	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo, Dir: "../../.."}
	pkgs, err := packages.Load(cfg, "./packages/tui/...")
	if err != nil {
		log.Fatalf("packages.Load: %v", err)
	}

	fields := map[string]map[string]types.Type{} // typeName -> field -> last type
	funcs := map[string]map[int]int{}            // funcName -> argIndex -> count seen (for arity)

	for _, pkg := range pkgs {
		for i, f := range pkg.Syntax {
			file := pkg.Fset.File(f.Pos()).Name()
			_ = i
			ast.Inspect(f, func(n ast.Node) bool {
				sel, ok := n.(*ast.SelectorExpr)
				if ok {
					// sel.Sel is identifier; check if X refers to package opencode
					if id, ok := sel.X.(*ast.Ident); ok && id.Name == "opencode" {
						name := sel.Sel.Name
						// record as used symbol
						_ = name
					}
				}

				// look for selector expressions like var.Field
				if se, ok := n.(*ast.SelectorExpr); ok {
					typ := pkg.TypesInfo.TypeOf(se.X)
					if typ == nil {
						return true
					}
					named, ok := typ.(*types.Named)
					if !ok {
						// check pointer
						if ptr, ok := typ.(*types.Pointer); ok {
							if nn, ok := ptr.Elem().(*types.Named); ok {
								named = nn
								ok = true
							}
						}
					}
					if ok {
						pkgPath := named.Obj().Pkg()
						if pkgPath != nil && pkgPath.Name() == "opencode" {
							typeName := named.Obj().Name()
							field := se.Sel.Name
							m := fields[typeName]
							if m == nil {
								m = map[string]types.Type{}
								fields[typeName] = m
							}
							fType := pkg.TypesInfo.TypeOf(se)
							if fType != nil {
								m[field] = fType
							}
						}
					}
				}

				// function calls
				if call, ok := n.(*ast.CallExpr); ok {
					switch fun := call.Fun.(type) {
					case *ast.SelectorExpr:
						if id, ok := fun.X.(*ast.Ident); ok && id.Name == "opencode" {
							fname := fun.Sel.Name
							m := funcs[fname]
							if m == nil {
								m = map[int]int{}
								funcs[fname] = m
							}
							m[len(call.Args)]++
						}
					}
				}

				return true
			})
			_ = file
		}
	}

	// Build output to patch api.md
	outPath := filepath.Join("..", "api.md")
	orig, _ := os.ReadFile(outPath)
	buf := bytes.NewBuffer(orig)
	buf.WriteString("\n\n<!-- INFERRED FIELDS -->\n")

	// sort types
	typesList := make([]string, 0, len(fields))
	for k := range fields {
		typesList = append(typesList, k)
	}
	sort.Strings(typesList)
	for _, t := range typesList {
		fmt.Fprintf(buf, "\n### %s\n\n", t)
		buf.WriteString("```go\n")
		fmt.Fprintf(buf, "type %s struct {\n", t)
		fmap := fields[t]
		fnames := make([]string, 0, len(fmap))
		for fn := range fmap {
			fnames = append(fnames, fn)
		}
		sort.Strings(fnames)
		for _, fn := range fnames {
			ft := fmap[fn]
			var tname string
			if ft != nil {
				tname = ft.String()
			} else {
				tname = "any"
			}
			fmt.Fprintf(buf, "    %s %s\n", fn, tname)
		}
		buf.WriteString("}\n```")
	}

	// functions
	buf.WriteString("\n\n<!-- INFERRED FUNCTIONS -->\n")
	fnNames := make([]string, 0, len(funcs))
	for k := range funcs {
		fnNames = append(fnNames, k)
	}
	sort.Strings(fnNames)
	for _, fn := range fnNames {
		// pick most common arity
		m := funcs[fn]
		bestA := 0
		bestCnt := 0
		for a, c := range m {
			if c > bestCnt {
				bestA = a
				bestCnt = c
			}
		}
		fmt.Fprintf(buf, "\nfunc %s(%s) (%s) { /* stub */ }\n", fn, genParams(bestA), "any, error")
	}

	err = os.WriteFile(outPath, buf.Bytes(), 0644)
	if err != nil {
		log.Fatalf("write api.md: %v", err)
	}
	fmt.Println("api.md updated with inferred fields and functions")
}

func genParams(n int) string {
	if n == 0 {
		return ""
	}
	parts := make([]string, n)
	for i := range n {
		parts[i] = fmt.Sprintf("arg%d any", i)
	}
	return join(parts, ", ")
}

func join(a []string, sep string) string {
	if len(a) == 0 {
		return ""
	}
	b := a[0]
	for i := 1; i < len(a); i++ {
		b += sep + a[i]
	}
	return b
}
