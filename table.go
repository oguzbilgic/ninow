// Package ninow implements a magical sql database ORM.
package ninow

import (
	"database/sql"
	"reflect"
	"strings"
)

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
	err := row.Scan(fieldPointers(t.fields, value)...)

	return value.Interface(), err
}

func (t *Table) SelectAllBy(column, value string) (interface{}, error) {
	query := "SELECT " + csv(t.columns) + " FROM " + t.name

	if column != "" && value != "" {
		// TODO Make sure this line is safe
		query += " WHERE " + column + "=" + value
	}

	rows, err := t.db.Query(query)
	if err != nil {
		return nil, err
	}

	sliceOfValue := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(t.rtype)), 0, 10)

	for rows.Next() {
		value := reflect.New(t.rtype)
		err := rows.Scan(fieldPointers(t.fields, value)...)
		if err != nil {
			return nil, err
		}

		sliceOfValue = reflect.Append(sliceOfValue, value)
	}

	return sliceOfValue.Interface(), nil
}

func (t *Table) Delete(id int) error {
	query := "DELETE FROM " + t.name + " WHERE id =? "

	_, err := t.db.Exec(query, id)
	return err
}
