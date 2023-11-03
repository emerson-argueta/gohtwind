package infra

import (
	"fmt"
	"html/template"
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

func formFunc(model interface{}, action string) template.HTML {
	// creates a html form based on the model
	// the fields of the model must be tagged with "form"
	// the action is the url to post the form to
	// the form is returned as a template.HTML
	// do code now
	modelType := reflect.TypeOf(model).Field(0).Type
	modelValue := reflect.ValueOf(model).Field(0)
	form := "<form action=\"" + action + "\" method=\"POST\">"
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fmt.Println("model value:", field)
		fmt.Println("field:", field)
		tag := field.Tag.Get("form")
		if tag == "-" {
			continue
		}
		value := getValue(modelValue.Field(i))
		form += "<label>" + tag + "</label>"
		form += "<input type=\"text\" name=\"" + tag + "\" value=\"" + value + "\">"
	}
	form += "<input type=\"submit\" value=\"Submit\">"
	form += "</form>"
	return template.HTML(form)

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
