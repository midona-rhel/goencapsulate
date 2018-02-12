package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "goencapsulate"
	app.Usage = "encapsulate third party packages"
	app.Flags = buildOutputFlags
	app.Action = func(c *cli.Context) error {
		s := c.String("t")
		if s == "" {
			println("provide t")
		} else {
			parse(s, []string{"Session"})
		}
		return nil
	}

	app.Run(os.Args)
}

func parse(s string, types []string) {
	fset := token.NewFileSet()
	dir, err := parser.ParseDir(fset, s, nil, 0)
	if err != nil {
		panic(err)
	}

	mergeFiles(dir)
	// var p *ast.Package

	// for _, pkg := range dir {
	// 	if !strings.HasSuffix(pkg.Name, "_test") {
	// 		p = pkg
	// 	}
	// }
	// ast.PackageExports(p)
	// ast.FilterPackage(p, func(s string) bool {
	// 	for _, t := range types {
	// 		if t == s {
	// 			return true
	// 		}
	// 	}
	// 	return false
	// })
	// for _, f := range p.Files {
	// 	var buf bytes.Buffer
	// 	if err := format.Node(&buf, fset, f); err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Printf("%s", buf.Bytes())
	// }
}

func cliTask(c *cli.Context) error {
	return nil
}

var buildOutputFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "target, t",
		Usage: "the target folder with the go source code",
	},

	cli.StringFlag{
		Name:  "output, o",
		Usage: "the output filename",
	},

	cli.StringSliceFlag{
		Name:  "export, e",
		Usage: "the types to be exported from the target source",
	},
}

func mergeFiles(dir map[string]*ast.Package) {
	// var nodes *ast.Node

	for _, pkg := range dir {
		if !strings.HasSuffix(pkg.Name, "_test") {
			ast.PackageExports(pkg)
			for _, f := range pkg.Files {
				ast.Inspect(f, func(node ast.Node) bool {
					if f, ok := node.(*ast.FuncDecl); ok {
						if f.Recv != nil && f.Recv.NumFields() == 1 {
							if name, err := typeName(f.Recv.List[0].Type); err == nil && name == "Session" {
								fmt.Println(f.Name.Name, name)
							}
						}
						return true
					}

					if t, ok := node.(*ast.TypeSpec); ok && t.Name.Name == "Session" {
						fmt.Println(t.Name.Name)
					}
					return true
				})
			}
		}
	}
}

// typeName returns the name of the type referenced by typeExpr.
func typeName(typeExpr ast.Expr) (string, error) {
	switch typeExpr := typeExpr.(type) {
	case *ast.StarExpr:
		return typeName(typeExpr.X)
	case *ast.Ident:
		return typeExpr.Name, nil
	default:
		return "", fmt.Errorf("expr %+v is not a type expression", typeExpr)
	}
}
