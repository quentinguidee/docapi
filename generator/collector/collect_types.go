package collector

import (
	"docapi/types"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type TypesCollector struct {
	// Structs are all the structs found in the project.
	Structs map[string]types.Value
}

func NewTypesCollector() *TypesCollector {
	return &TypesCollector{
		Structs: map[string]types.Value{},
	}
}

func (a *TypesCollector) Run(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		return a.collect(path)
	})
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
			case *ast.StructType:
				id := x.Name.Name
				a.Structs[id] = types.Value{
					Type:   "object",
					Fields: map[string]types.Value{},
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
							t = "array"
						case *ast.MapType:
							t = "object"
						default:
							t = "unknown"
						}

						a.Structs[id].Fields[jsonName] = types.Value{
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

func (a *TypesCollector) Output() (map[string]types.Value, error) {
	return a.Structs, nil
}

func isDefaultType(name string) bool {
	switch name {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
		"float32", "float64", "complex64", "complex128",
		"string", "bool", "byte", "rune":
		return true
	default:
		return false
	}
}
