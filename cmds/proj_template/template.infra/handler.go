package infra

import (
	"net/http"
	"reflect"
	"strconv"
	"time"
)

func UnmarshalForm(r *http.Request, dst interface{}) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	val := reflect.ValueOf(dst).Elem()
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		fieldName := typ.Field(i).Tag.Get("form")
		if fieldName == "-" {
			continue
		}
		formVal := r.FormValue(fieldName)
		if formVal == "" {
			continue
		}
		field := val.Field(i)
		err = setField(field, formVal)
		if err != nil {
			return err
		}
	}
	return nil
}

func setField(field reflect.Value, formVal string) error {
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
	case reflect.Struct:
		if field.Type().String() == "time.Time" {
			timeVal := field.Interface().(time.Time)
			timeVal, err := time.Parse("2006-01-02T15:04", formVal)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(timeVal))
		}
	}
	return nil
}
