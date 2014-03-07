package ninow

import (
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

func csQ(n int) string {
	questions := make([]string, n)
	for i := 0; i < n; i++ {
		questions[i] = "?"
	}
	return csv(questions)
}

func cscv(columns []string, values []interface{}) string {
	queryString := ""

	for i, column := range columns {
		queryString += column + "=?"
		if i < len(columns)-1 {
			queryString += ", "
		}
	}

	return queryString
}

func fieldValues(fields []string, value reflect.Value) []interface{} {
	indirectValue := reflect.Indirect(value)

	var values []interface{}
	for _, field := range fields {
		fieldValue := indirectValue.FieldByName(field).Interface()
		values = append(values, fieldValue)
	}

	return values
}

func fieldPointers(fields []string, value reflect.Value) []interface{} {
	indirectValue := reflect.Indirect(value)

	var pointers []interface{}
	for _, field := range fields {
		fieldValue := indirectValue.FieldByName(field)
		addr := fieldValue.Addr()
		pointers = append(pointers, addr.Interface())
	}

	return pointers
}

func fieldNameToColumnName(fieldName string) string {
	//Fix edge cases
	for _, keyword := range []string{"ID", "URL"} {
		if !strings.Contains(fieldName, keyword) {
			continue
		}

		weedle := strings.Index(fieldName, keyword)
		keywordLength := len(keyword)
		fieldName = fieldName[0:weedle+1] + strings.ToLower(fieldName[weedle+1:weedle+keywordLength]) + fieldName[weedle+keywordLength:len(fieldName)]
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
