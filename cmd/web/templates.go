package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	libw "github.com/alicebob/w/w"
)

var (
	baseTempl = template.Must(
		template.New("base").
			Funcs(template.FuncMap{
				"title": libw.Title,
				"link": func(s string) string {
					return *baseURL + s
				},
			}).Parse(`<!DOCTYPE html>
<html>
    <head>
        <title>{{ .title }}</title>
        {{- template "head" . }}
    </head>
    <body>
		<a href="{{link "/"}}">Home</a><br />
		<hr />
        {{- template "page" . }}
    </body>
</html>
{{define "head"}}
{{- end}}
{{define "page"}}
{{- end}}
`))

	indexTempl = template.Must(extend(baseTempl).Parse(`
{{define "page"}}
Hello World<br />

<br />
<a href="./adhoc/">ad hoc feeds</a><br />
<a href="./curated/">curated feeds</a><br />
<br />

Recent new versions:<br />
	{{- range .entries}}
		<a href="./p/{{.Page}}/">{{title .Page}}</a>: {{.StableVersion}}<br />
	{{- end}}
{{- end}}
`))
)

func extend(t *template.Template) *template.Template {
	return template.Must(t.Clone())
}

func runTmpl(w http.ResponseWriter, t *template.Template, args interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	b := &bytes.Buffer{}
	if err := t.Execute(b, args); err != nil {
		log.Printf("template: %s", err)
		http.Error(w, "internal server error", 500)
		return
	}
	w.Write(b.Bytes())
}
