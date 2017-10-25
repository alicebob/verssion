package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
)

var (
	baseTempl = template.Must(template.New("base").Parse(`
<!DOCTYPE html>
<html>
    <head>
        <title>{{ .Title }}</title>
        <link rel="stylesheet" href="/static/s.css" />
    </head>
    <body>
        {{- template "page" . }}
    </body>
</html>
{{define "page"}}
{{- end}}
`))

	indexTempl = template.Must(extend(baseTempl).Parse(`
{{define "page"}}
Hello World<br />
	{{- range .entries}}
		<a href="./v/{{.Page}}">{{.Page}}</a>: {{.StableVersion}}<br />
	{{- end}}
{{- end}}
`))

	pageTempl = template.Must(extend(baseTempl).Parse(`
{{define "page"}}
	{{.page}}<br />
	{{- range .versions}}
		{{- .StableVersion}} - (spider: {{.T}})<br />
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
