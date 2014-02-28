// Package ninow implements a magical sql database ORM.
package ninow

import (
	"database/sql"
	"reflect"
	"strings"
	"unicode"
)

func csv(values []string) string {
	list := ""
	for i, value := range values {
		if i != (len(values) - 1) {
			list += value + ", "
		} else {
			list += value
		}
	}
	return list
}

type Table struct {
	db    *sql.DB
	rtype reflect.Type

	name    string
	columns []string
	fields  []string
}

func TableFor(row interface{}, db *sql.DB) *Table {
	rtype := reflect.TypeOf(row)
	name := strings.ToLower(rtype.Name()) + "s"

	columns := []string{}
	fields := []string{}
	for i := 0; i < rtype.NumField(); i++ {
		fieldValue := rtype.Field(i)
		//TODO fix this piece
		columns = append(columns, fieldNameToColumnName(fieldValue.Name))
		fields = append(fields, fieldValue.Name)
	}

	return &Table{db, rtype, name, columns, fields}
}

func (t *Table) Select(id int) (interface{}, error) {
	row := t.db.QueryRow("SELECT "+csv(t.columns)+" FROM "+t.name+" WHERE id=?", id)
	value := reflect.New(t.rtype)
	indirectValue := reflect.Indirect(value)

	var pointers []interface{}
	for _, field := range t.fields {
		fieldValue := indirectValue.FieldByName(field)
		addr := fieldValue.Addr()
		pointers = append(pointers, addr.Interface())
	}

	err := row.Scan(pointers...)
	if err != nil {
		return nil, err
	}

	return value.Interface(), nil
}

func fieldNameToColumnName(fieldName string) string {
	//TODO add more edge cases to conform idiomatically
	edgeCases := []string{"ID", "URL"}

	for _, keyword := range edgeCases {
		fieldName = fixEdgeFieldName(fieldName, keyword)
	}

	columnName := strings.ToLower(string(fieldName[0]))
	for _, char := range []rune(fieldName[1:]) {
		if unicode.IsUpper(char) {
			columnName = columnName + "_" + strings.ToLower(string(char))
		} else {
			columnName = columnName + string(char)
		}
	}
	return columnName
}

func fixEdgeFieldName(fieldName string, keyword string) string {
	if !strings.Contains(fieldName, keyword) {
		return fieldName
	}
	weedle := strings.Index(fieldName, keyword)
	keywordLength := len(keyword)
	return fieldName[0:weedle+1] + strings.ToLower(fieldName[weedle+1:weedle+keywordLength]) + fieldName[weedle+keywordLength:len(fieldName)]

}
