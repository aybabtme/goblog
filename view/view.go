package view

import (
	"fmt"
	"html/template"
)

func getApplicationTemplate() (*template.Template, error) {
	app, err := template.ParseFiles("view/template/application.tmpl")
	if nil != err {
		return nil, err
	}

	base, err := app.ParseGlob("view/template/base/*.tmpl")
	if nil != err {
		return nil, err
	}
	return base, nil
}

func getTemplateByName(templateName string) (*template.Template, error) {
	base, err := getApplicationTemplate()
	if nil != err {
		fmt.Println("Cannot prepare application template.", err)
		return nil, err
	}

	content, err := base.ParseFiles(
		fmt.Sprintf("view/template/%s.tmpl", templateName))
	if nil != err {
		fmt.Printf("Couldn't load %s template.\n %v", templateName, err)
		return nil, err
	}
	return content, nil
}

func GetIndexTemplate() *template.Template {
	return template.Must(getTemplateByName("index"))
}

func GetAuthorTemplate() *template.Template {
	return template.Must(getTemplateByName("author"))
}

func GetPostListingTemplate() *template.Template {
	return template.Must(getTemplateByName("post_listing"))
}

func GetPostTemplate() *template.Template {
	return template.Must(getTemplateByName("post"))
}

func GetLabelTemplate() *template.Template {
	return template.Must(getTemplateByName("label"))
}

func GetUserTemplate() *template.Template {
	return template.Must(getTemplateByName("user"))
}
