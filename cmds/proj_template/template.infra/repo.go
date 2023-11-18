package infra

import (
	mysql "github.com/go-jet/jet/v2/mysql"
	"reflect"
	"time"
)

type Pagination struct {
	TotalRecords int64 `json:"total_records"`
	TotalPages   int64 `json:"total_pages"`
	CurrentPage  int64 `json:"current_page"`
	PerPage      int64 `json:"per_page"`
}

func NewPagination(current_page int64, per_page int64, total_records int64) *Pagination {
	total_pages := total_records / per_page
	return &Pagination{
		TotalRecords: total_records,
		TotalPages:   total_pages,
		CurrentPage:  current_page,
		PerPage:      per_page,
	}
}

func UpdateColumns(model interface{}) mysql.ColumnList {
	v := reflect.ValueOf(model)
	t := v.Type()

	// If v is a pointer, dereference it to get the underlying value
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}

	var updateSet mysql.ColumnList
	for i := 0; i < v.NumField(); i++ {
		fieldType := t.Field(i)

		// Skip unexported fields and fields tagged with form:"-"
		if fieldType.PkgPath != "" || fieldType.Tag.Get("form") == "-" || fieldType.Tag.Get("form") == "" {
			continue
		}

		// Assume column names are the same as field names
		columnName := fieldType.Name

		// Check the type of the field and add the appropriate type of column
		switch fieldType.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			updateSet = append(updateSet, mysql.IntegerColumn(columnName))
		case reflect.Float32, reflect.Float64:
			updateSet = append(updateSet, mysql.FloatColumn(columnName))
		case reflect.String:
			updateSet = append(updateSet, mysql.StringColumn(columnName))
		case reflect.Struct:
			if fieldType.Type == reflect.TypeOf(time.Time{}) {
				updateSet = append(updateSet, mysql.TimestampColumn(columnName))
			}
		}
	}

	return updateSet
}
