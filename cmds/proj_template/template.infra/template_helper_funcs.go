package infra

import (
	"fmt"
	"html/template"
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

var TemplateHelperFuncs = template.FuncMap{"dict": dictFunc, "slice": sliceFunc}
