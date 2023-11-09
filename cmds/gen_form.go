package cmds

import (
	"flag"
	"fmt"
	"gohtwind/utils"
	"os"
	"strings"
)

func genFormUsageString() string {
	return `
Usage: gohtwind gen-form [options]
	Options:
		-feature-name string (required)
			Name of the feature the form is for
		-model-name string (required)
			Name of the model the form is for
		-template-name string (required)
			Name of the template the form is for
		-instance-name string (optional)
			Name of the variable holding the instance of the model in the template (useful for update forms)
			omit this option for create forms
		-action string (required)
			Form action attribute
	Info:
		* Replaces {{GEN_FORM}} in the specified template with a form for the specified model
		* The form is generated using the model's form tags
		* When an instance name is provided, the form is generated with the instance's values 
`
}

type form struct {
	flagSet      *flag.FlagSet
	featName     *string
	modelName    *string
	templateName *string
	instanceName *string
	action       *string
	projectPath  string
}

func newForm() *form {
	genFormFlags := flag.NewFlagSet("gohtwind gen-form", flag.ExitOnError)
	pjp, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	featName := genFormFlags.String("feature-name", "", "Name of the feature the form is for")
	modelName := genFormFlags.String("model-name", "", "Name of the model the form is for")
	templateName := genFormFlags.String("template-name", "", "Name of the template the form is for")
	instanceName := genFormFlags.String("instance-name", "", "Name of the variable holding the instance of the model in the template (useful for update forms)")
	action := genFormFlags.String("action", "", "Form action attribute")
	return &form{
		flagSet:      genFormFlags,
		featName:     featName,
		modelName:    modelName,
		templateName: templateName,
		instanceName: instanceName,
		action:       action,
		projectPath:  pjp,
	}
}

func GenForm() {
	f := newForm()
	args := os.Args[2:]
	err := f.flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if *f.featName == "" || *f.modelName == "" || *f.templateName == "" || *f.action == "" {
		fmt.Println(genFormUsageString())
		os.Exit(1)
	}

	f.genForm()
}

func (f *form) genForm() {
	tp := fmt.Sprintf("%s/%s/templates/%s", f.projectPath, *f.featName, *f.templateName)
	tf, err := os.ReadFile(tp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	content := string(tf)
	content = strings.ReplaceAll(content, "{{GEN_FORM}}", f.genFormString())
	err = os.WriteFile(tp, []byte(content), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (f *form) genFormString() string {
	fp := fmt.Sprintf("%s/%s/dtos.go", f.projectPath, *f.featName)
	structFields, err := utils.GetStructFields(fp, *f.modelName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if *f.instanceName != "" {
		res := fmt.Sprintf("<form action=\"%s/{{.%s.ID}}\" method=\"PATCH\">\n", *f.action, *f.instanceName)
		res += f.genFormFieldsWithInstance(structFields)
		return res
	}

	res := fmt.Sprintf("<form action=\"%s\" method=\"POST\">\n", *f.action)
	res += f.genFormFields(structFields)
	return res
}

func (f *form) genFormFieldsWithInstance(structFields []utils.FieldInfo) string {
	res := fmt.Sprintf("\t<input type=\"hidden\" name=\"_method\" value=\"PATCH\">\n")
	for _, field := range structFields {
		form_tag := findFormTag(field.Tag)
		typ := determineFormFieldType(field.Type)
		res += fmt.Sprintf("\t<label>%s</label>\n", form_tag.Value)
		res += fmt.Sprintf("\t<input type=\"%s\" name=\"%s\" value=\"{{.%s.%s}}\">\n", typ, form_tag.Value, *f.instanceName, field.Name)
	}
	res += "\t<input type=\"submit\" value=\"Submit\">\n"
	res += "</form>\n"
	return res
}

func (f *form) genFormFields(structFields []utils.FieldInfo) string {
	res := "\t<input type=\"hidden\" name=\"_method\" value=\"POST\">\n"
	for _, field := range structFields {
		form_tag := findFormTag(field.Tag)
		typ := determineFormFieldType(field.Type)
		res += fmt.Sprintf("\t<label>%s</label>\n", form_tag.Value)
		res += fmt.Sprintf("\t<input type=\"%s\" name=\"%s\">\n", typ, form_tag.Value)
	}
	res += "\t<input type=\"submit\" value=\"Submit\">\n"
	res += "</form>\n"
	return res
}

func findFormTag(tags []struct{ Name, Value string }) struct{ Name, Value string } {
	for _, tag := range tags {
		if tag.Name == "form" {
			return tag
		}
	}
	return struct{ Name, Value string }{}
}

func determineFormFieldType(fieldType string) string {
	switch fieldType {
	case "string":
		return "text"
	case "int64":
		return "number"
	case "float64":
		return "number"
	case "time.Time":
		return "datetime-local"
	default:
		return "text"
	}
}
