package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	os "os"
	"reflect"
	"strconv"
	"strings"
)

type tableStruct struct {
	node      *ast.StructType
	tableName string
}

// Column ...
type Column struct {
	pk      string
	unique  string
	notnull string
	created string
	updated string
	name    string
	t       string
}

func isGenTableDoc(doc *ast.CommentGroup) (string, bool) {
	if doc == nil {
		return "", false
	}

	for _, comment := range doc.List {
		if comment == nil {
			continue
		}

		value := strings.TrimSpace(comment.Text)
		if strings.HasPrefix(value, "//") {
			value = strings.TrimSpace(value[2:])
			if strings.HasPrefix(value, "+genTable:") {
				return strings.TrimSpace(value[10:]), true
			}
		} else {
			// multi line
			lines := strings.Split(value, "\n")
			for _, line := range lines {
				value := strings.TrimSpace(line)
				if strings.HasPrefix(value, "+genTable:") {
					return strings.TrimSpace(value[10:]), true
				}
			}

		}
	}

	return "", false
}

func filterGenDeclTableStruct(n *ast.GenDecl) ([]*tableStruct, error) {
	tableName, isGenTable := isGenTableDoc(n.Doc)
	tables := []*tableStruct{}

	for _, spec := range n.Specs {
		if n, ok := spec.(*ast.TypeSpec); ok {
			typeTableName, typeGenTable := isGenTableDoc(n.Doc)
			if !typeGenTable && !isGenTable {
				fmt.Println(n.Name, " doesn't have genTable spec, skip it")
				continue
			}
			if tableName != "" && typeTableName != "" {
				return nil, fmt.Errorf("%s has multi spec tableName", n.Name)
			}
			if tableName != "" && len(tables) > 0 {
				return nil, fmt.Errorf("%s spec multi struct", n.Name)
			}

			structType, ok := n.Type.(*ast.StructType)
			if !ok {
				if typeTableName != "" {
					return nil, fmt.Errorf("%s genTable spec on %v", n.Name, reflect.TypeOf(n.Type).Elem().Name())
				}
				fmt.Println(n.Name, " isn't struct, skip it")
				continue
			}

			if typeTableName == "" {
				typeTableName = tableName
			}
			if typeTableName == "" {
				typeTableName = n.Name.Name
			}

			tables = append(tables, &tableStruct{
				tableName: typeTableName,
				node:      structType,
			})
		}
	}

	if isGenTable && len(tables) == 0 {
		return nil, fmt.Errorf("has genTable spec but on table exists")
	}

	return tables, nil
}

func filterFileTableStruct(file *ast.File) ([]*tableStruct, error) {
	genTables := []*tableStruct{}

	for _, decl := range file.Decls {
		var n ast.Node = decl
		if n, ok := n.(*ast.GenDecl); ok {
			tables, err := filterGenDeclTableStruct(n)
			if err != nil {
				return nil, err
			}
			genTables = append(genTables, tables...)
		}
	}

	return genTables, nil
}

func genColumn(field *ast.Field) (*Column, error) {
	col := &Column{}

	if len(field.Names) > 1 {
		return nil, fmt.Errorf("%v  should only be one", field.Names)
	}

	if len(field.Names) == 1 {
		col.name = field.Names[0].Name
	} else {
		ident, ok := field.Type.(*ast.Ident)
		if !ok {
			return nil, fmt.Errorf("field doesn't have fieldNames, and type isn't Ident, skip it")
		}
		col.name = ident.Name
		col.t = col.name
	}

	tag := ""
	if field.Tag != nil {
		var err error
		tag, err = strconv.Unquote(field.Tag.Value)
		if err != nil {
			fmt.Println("unquote tag failed, %v", err)
			os.Exit(1)
		}
		structTag, ok := reflect.StructTag(tag).Lookup("xorm")
		if ok {
			tag = structTag
		} else {
			tag = ""
		}
	}

	tags := strings.Split(tag, " ")
	for _, tagField := range tags {
		if tagField == "" {
			continue
		}

		switch tagField {
		case "notnull":
			col.notnull = "NOT NULL"
		case "unique":
			col.unique = "UNIQUE KEY"
		case "pk":
			col.pk = "PRIMARY KEY"
		case "created":
			col.created = "DEFAULT CURRENT_TIMESTAMP"
		case "updated":
			col.updated = "ON UPDATE CURRENT_TIMESTAMP"
		case "int":
			col.t = "INT"
		case "bigint":
			col.t = "BIGINT"
		default:
			if strings.HasPrefix(tagField, "'") {
				if len(tagField) <= 2 {
					return nil, fmt.Errorf("%v used as fieldname but doesn't have value", tagField)
				}
				if !strings.HasSuffix(tagField, "'") {
					return nil, fmt.Errorf("%v used as fieldname but doesn't close quote", tagField)
				}
				tagField = tagField[1 : len(tagField)-1]
				col.name = tagField
			} else {
				if strings.HasPrefix(tagField, "varchar") {
					col.t = tagField
				} else {
					return nil, fmt.Errorf("unknown tag: %v", tagField)
				}
			}
		}
	}

	if col.t == "" {
		switch ident := field.Type.(type) {
		case *ast.Ident:
			switch ident.Name {
			case "int64":
				col.t = "BIGINT"
			case "time.Time":
				col.t = "DATETIME"
			case "string":
				col.t = "VARCHAR(255)"
			case "int":
				col.t = "INT"
			default:
				return nil, fmt.Errorf("%v type to sqltype failed, unknown: %v", col.name, ident.Name)
			}
		case *ast.SelectorExpr:
			if ident.Sel.Name == "Time" {
				if ident, ok := ident.X.(*ast.Ident); ok {
					if ident.Name == "time" {
						col.t = "DATETIME"
					}
				}
			}
		default:
			return nil, fmt.Errorf("%v field type unknown: %v", col.name, reflect.TypeOf(field.Type).Elem().Name())
		}
	}

	if col.t == "" {
		return nil, fmt.Errorf("%v doesn't have type", col.name)
	}

	return col, nil
}

func genColumns(table *tableStruct) ([]*Column, error) {
	columns := []*Column{}

	for _, field := range table.node.Fields.List {
		col, err := genColumn(field)
		if err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	for _, col := range columns {
		if col.name == "Id" && col.t == "BIGINT" {
			col.pk = "PRIMARY KEY"
		}
	}

	if 0 == len(columns) {
		return nil, fmt.Errorf("no column")
	}

	return columns, nil
}

func genSql(table *tableStruct) (string, error) {
	columns, err := genColumns(table)
	if err != nil {
		return "", err
	}

	// 开始生成sql
	builder := strings.Builder{}
	builder.WriteString("CREATE TABLE `")
	builder.WriteString(table.tableName)
	builder.WriteString("` (")

	for i, col := range columns {
		builder.WriteString("\n\t`")
		builder.WriteString(col.name)
		builder.WriteString("` ")
		builder.WriteString(col.t)
		if col.notnull != "" {
			builder.WriteString(" NOT NULL")
		}
		if col.created != "" {
			builder.WriteString(" " + col.created)
		}
		if col.updated != "" {
			builder.WriteString(" " + col.updated)
		}

		if i != len(columns)-1 {
			builder.WriteString(",")
		}
	}

	for _, col := range columns {
		if col.pk != "" {
			builder.WriteString(",\n\t")
			builder.WriteString("PRIMARY KEY(`")
			builder.WriteString(col.name)
			builder.WriteString("`)")
		}
		if col.unique != "" {
			builder.WriteString(",\n\t")
			builder.WriteString("UNIQUE KEY `")
			builder.WriteString(col.name)
			builder.WriteString("` (`")
			builder.WriteString(col.name)
			builder.WriteString("`)")
		}
	}

	builder.WriteString("\n);")
	return builder.String(), nil
}

func createTable(file string, src interface{}) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, file, src, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	genTables, err := filterFileTableStruct(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, table := range genTables {
		sql, err := genSql(table)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(sql)
	}
}
