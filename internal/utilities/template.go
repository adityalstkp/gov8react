package utilities

import "text/template"

func CreateTemplate(path string, templateName string) (*template.Template, error) {
	fByte, err := ReadFile(path)
	if err != nil {
		return nil, err
	}

	tmpl := template.Must(template.New(templateName).Parse(string(fByte)))
	return tmpl, nil

}
