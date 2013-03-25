package view

import (
	"html/template"
	"os"
	"testing"
)

func TestApplicationTemplate(t *testing.T) {

	app, err := template.ParseFiles("template/application.tmpl")
	if nil != err {
		t.Error("Couldn't load templates.", err)
		return
	}

	base, err := app.ParseGlob("template/base/*.tmpl")
	if nil != err {
		t.Error("Couldn't load templates.", err)
		return
	}

	content, err := base.Parse("{{define \"content\"}}<h1>Hello World!</h1>{{end}}")
	if nil != err {
		t.Error("Couldn't load templates.", err)
		return
	}

	if err := content.Execute(os.Stdout, nil); nil != err {
		t.Errorf("template execution: %s", err)
	}
}
