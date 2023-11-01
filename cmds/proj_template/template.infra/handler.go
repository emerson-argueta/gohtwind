package infra

import (
	"net/http"
	"reflect"
	"strconv"
)

func UnmarshalForm(r *http.Request, dst interface{}) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	val := reflect.ValueOf(dst).Elem()
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("form")
		if tag == "-" {
			continue
		}

		formVal := r.FormValue(tag)
		if formVal == "" {
			continue
		}

		field := val.Field(i)
		switch field.Kind() {
		case reflect.String:
			field.SetString(formVal)
		case reflect.Int64:
			intVal, err := strconv.ParseInt(formVal, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(intVal)
		case reflect.Float64:
			floatVal, err := strconv.ParseFloat(formVal, 64)
			if err != nil {
				return err
			}
			field.SetFloat(floatVal)
		}
	}
	return nil
}
