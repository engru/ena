// sfsdfsdfsf package

// package main ...
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"reflect"

	// xxxx
	os "os"
)

// 注释

func printDoc(doc *ast.CommentGroup) {
	if doc != nil {
		for _, comment := range doc.List {
			if comment != nil {
				fmt.Println("  ", comment.Text)
			}
		}
	}
}

// FirstType docs
type FirstType struct { // first type docs inline
	// FirstMember docs
	FirstMember string `json:""`
}

// SecondType docs
type SecondType struct {
	// SecondMember docs
	SecondMember, member string // line comment
}

func parseMain() {
	data, err := ioutil.ReadFile("./main.go")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// comment2
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "main.go", data, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// ast.Print(fs, f)
	ast.Inspect(f, func(n ast.Node) bool {
		fmt.Println("---------------")
		if n != nil {
			fmt.Println("Reflect Type: ", reflect.TypeOf(n).Elem().Name())
		}

		switch x := n.(type) {
		case *ast.FuncDecl:
			fmt.Println("FuncDecl: ", x.Name, "    ")
			printDoc(x.Doc)
			return true
		case *ast.ImportSpec:
			fmt.Println("ImportSpec: ", x.Name, "  ", x.Path)
			printDoc(x.Doc)
			return true
		case *ast.Comment:
			fmt.Println("Comment: ", x.Text)
		case *ast.File:
			fmt.Println("File: ", x.Name)
			printDoc(x.Doc)
			return true
		case *ast.TypeSpec:
			fmt.Println("Struct: ", x.Name, "  ", reflect.TypeOf(x.Type).Elem().Name())
			printDoc(x.Doc)
			printDoc(x.Comment)
		case *ast.StructType:
			fmt.Println("StructType: ")
		case *ast.Field:
			fmt.Println("Field: ", x.Names)
			fmt.Printf("Tag: %+v\n", x.Tag)
			printDoc(x.Doc)
			printDoc(x.Comment)
		case *ast.GenDecl:
			fmt.Println("GenDecl: ")
			printDoc(x.Doc)
			return true
		case *ast.CommentGroup:
			fmt.Println("CommentGroup: ")
		case ast.Node:
			fmt.Println("Node: ")
		}
		return true
	})

	for _, c := range f.Comments {
		printDoc(c)
	}
}

// comment2
// sdfsdf
//  werwer
func main() {
	createTable("./table.go", nil)
}
