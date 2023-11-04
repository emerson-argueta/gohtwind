package infra

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"reflect"
	"time"
)

func dictFunc(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func sliceFunc(values ...interface{}) []interface{} {
	return values
}

func formFunc(model interface{}, action string, method string) template.HTML {
	modelType := reflect.TypeOf(model).Field(0).Type
	modelValue := reflect.ValueOf(model).Field(0)
	form := fmt.Sprintf("<form action=\"%s\" method=\"%s\">", action, method)
	if method == "PATCH" {
		form += fmt.Sprintf(`<input type="hidden" name="_method" value="%s">`, method)
	}
	tk, err := csrfToken()
	if err != nil {
		log.Printf("error generating csrf token: %v", err)
	}
	form += fmt.Sprintf(`<input type="hidden" name="csrf_token" value="%s">`, tk)
	for i := 0; i < modelType.NumField(); i++ {
		name := modelType.Field(i).Name
		value := getValue(modelValue.Field(i))
		form += fmt.Sprintf("<label>%s</label>", name)
		form += fmt.Sprintf("<input type=\"text\" name=\"%s\" value=\"%s\">", name, value)
	}
	form += "<input type=\"submit\" value=\"Submit\">"
	form += "</form>"
	return template.HTML(form)

}

func csrfToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func getValue(value reflect.Value) string {
	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Int64:
		return fmt.Sprintf("%d", value.Int())
	case reflect.Float64:
		return fmt.Sprintf("%f", value.Float())
	case reflect.Struct:
		if value.Type().String() == "time.Time" {
			return value.Interface().(time.Time).Format("2006-01-02T15:04")
		}
	}
	return ""

}

var TemplateHelperFuncs = template.FuncMap{
	"dict":  dictFunc,
	"slice": sliceFunc,
	"form":  formFunc,
}
