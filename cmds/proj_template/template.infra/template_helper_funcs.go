package infra

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"reflect"
	"strings"
	"time"
)

func unescapeJSFunc(s string) template.JS {
	return template.JS(s)
}

func unescapeHTMLFunc(s string) template.HTML {
	return template.HTML(s)
}

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
	modelType := reflect.TypeOf(model)
	modelValue := reflect.ValueOf(model)
	tpl := fmt.Sprintf(`
		<section>
		  <div class="py-8 px-4 mx-auto max-w-5xl lg:py-16">
			  <h2 class="mb-4 text-xl font-bold text-gray-900">Add a new %s</h2>
			  {{FORM_STR}}
		  </div>
		</section>
	`, modelType.Name())
	form := fmt.Sprintf("<form action=\"%s\" method=\"%s\">", action, method)
	form += `
            <div class="grid gap-4 sm:grid-cols-2 sm:gap-6">
	`
	if method == "PATCH" {
		form += fmt.Sprintf(`<input type="hidden" name="_method" value="%s">`, method)
	}
	tk, err := csrfToken()
	if err != nil {
		log.Printf("error generating csrf token: %v", err)
	}
	form += fmt.Sprintf(`<input type="hidden" name="csrf_token" value="%s">`, tk)
	vis_count := 0
	for i := 0; i < modelType.NumField(); i++ {
		name := modelType.Field(i).Tag.Get("form")
		value := getValue(modelValue.Field(i))
		// if the name is "-" then skip this field
		if name == "-" {
			continue
		}
		// if the name is id add as a hidden field
		if name == "id" {
			form += fmt.Sprintf(`<input type="hidden" name="%s" value="%s">`, name, value)
			continue
		}
		if vis_count%3 != 0 {
			form += `<div>`
			form += fmt.Sprintf("<label class=\"form-label\">%s</label>", name)
			form += fmt.Sprintf("<input class=\"form-input\" type=\"text\" name=\"%s\" value=\"%s\">", name, value)
			form += `</div>`
			vis_count++
			continue
		}
		form += `<div class="sm:col-span-2">`
		form += fmt.Sprintf("<label class=\"form-label\">%s</label>", name)
		form += fmt.Sprintf("<input class=\"form-input\" type=\"text\" name=\"%s\" value=\"%s\">", name, value)
		form += `</div>`
		vis_count++
	}
	form += `</div>`
	form += fmt.Sprintf("<button type=\"submit\" class=\"mt-4 sm:mt-6 btn-blue\">Add %s", modelType.Name())
	form += "</form>"
	tpl = strings.ReplaceAll(tpl, "{{FORM_STR}}", form)
	return template.HTML(tpl)

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
	"dict":         dictFunc,
	"slice":        sliceFunc,
	"form":         formFunc,
	"unescapeJS":   unescapeJSFunc,
	"unescapeHTML": unescapeHTMLFunc,
}
