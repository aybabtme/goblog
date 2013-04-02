package view

import (
	"fmt"
	"text/template"
)

/*
 * Helpers
 */
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

/*
 * Index
 */

func GetIndexTemplate() *template.Template {
	return template.Must(getTemplateByName("index"))
}

/*
 *	Posts
 */

func GetPostListingTemplate() *template.Template {
	return template.Must(getTemplateByName("post_listing"))
}

func GetPostTemplate() *template.Template {
	return template.Must(getTemplateByName("post"))
}

func GetPostComposeTemplate() *template.Template {
	return template.Must(getTemplateByName("post_compose"))
}

func GetPostDestroyTemplate() *template.Template {
	return template.Must(getTemplateByName("post"))
}

/*
 * Labels
 */

func GetLabelTemplate() *template.Template {
	return template.Must(getTemplateByName("label"))
}

/*
 * Users
 */

func GetUserTemplate() *template.Template {
	return template.Must(getTemplateByName("user"))
}

/*
 * Authors
 */

func GetAuthorTemplate() *template.Template {
	return template.Must(getTemplateByName("author"))
}

func GetAuthorListTemplate() *template.Template {
	return template.Must(getTemplateByName("author_listing"))
}
