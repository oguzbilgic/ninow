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
