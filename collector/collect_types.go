package collector

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Struct struct {
	Type   string
	Fields map[string]Struct
}

type TypesCollector struct {
	// Structs are all the structs found in the project.
	Structs map[string]Struct
	// Aliases are all the aliases found in the project.
	// e.g. type MyString string
	Aliases map[string]string
}

func NewTypesCollector() *TypesCollector {
	return &TypesCollector{
		Structs: map[string]Struct{},
		Aliases: map[string]string{},
	}
}

func (a *TypesCollector) Run(path string) (map[string]Struct, map[string]string, error) {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		return a.collect(path)
	})
	if err != nil {
		return nil, nil, err
	}
	return a.Structs, a.Aliases, nil
}

func (a *TypesCollector) collect(path string) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return err
	}
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			switch x.Type.(type) {
			case *ast.Ident:
				a.Aliases[x.Name.Name] = x.Type.(*ast.Ident).Name
			case *ast.StructType:
				id := x.Name.Name
				a.Structs[id] = Struct{
					Type:   "object",
					Fields: map[string]Struct{},
				}
				for _, field := range x.Type.(*ast.StructType).Fields.List {
					tag := field.Tag
					if tag == nil {
						continue
					}
					tags := strings.Split(tag.Value, " ")
					for _, tag := range tags {
						if !strings.HasPrefix(tag, "`json:") {
							continue
						}
						jsonName := strings.TrimPrefix(tag, "`json:")
						jsonName = strings.TrimSuffix(jsonName, "`")
						jsonName, err = strconv.Unquote(jsonName)
						jsonName = strings.Split(jsonName, ",")[0]
						if err != nil {
							log.Fatal(err)
						}
						if jsonName == "-" {
							continue
						}

						tp := field.Type
						if _, ok := field.Type.(*ast.StarExpr); ok {
							tp = field.Type.(*ast.StarExpr).X
						}

						var t string
						switch tp.(type) {
						case *ast.BasicLit:
							t = tp.(*ast.BasicLit).Value
						case *ast.SelectorExpr:
							t = tp.(*ast.SelectorExpr).Sel.Name
						case *ast.Ident:
							t = tp.(*ast.Ident).Name
						case *ast.ArrayType:
							id := tp.(*ast.ArrayType).Elt.(*ast.Ident).Name
							t = "[]" + id
						case *ast.MapType:
							t = "object"
						default:
							t = "unknown"
						}

						a.Structs[id].Fields[jsonName] = Struct{
							Type: t,
						}
					}
				}
			}
		}
		return true
	})
	return nil
}

func (a *TypesCollector) Output() (map[string]Struct, error) {
	return a.Structs, nil
}
