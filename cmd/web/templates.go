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
			}).Parse(`<!DOCTYPE html>
<html>
    <head>
        <title>{{ .title }}</title>
        {{- template "head" . }}
    </head>
    <body>
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
<br />

Recent new versions:<br />
	{{- range .entries}}
		<a href="./v/{{.Page}}/">{{title .Page}}</a>: {{.StableVersion}}<br />
	{{- end}}
{{- end}}
`))

	adhocTempl = template.Must(extend(baseTempl).Parse(`
{{define "head"}}
	<link rel="alternate" type="application/atom+xml" title="Atom 1.0" href="{{.atom}}"/>
{{- end}}
{{define "page"}}

{{with .errors}}
	&ltblink>Errors&lt/blink><br />
	{{- range .}}
		{{.}}<br />
	{{- end}}
	<br />
{{end}}

<form method="GET">
	Create an ad hoc Atom URL to track versions.<br />
	<br />
	<br />
	{{- if .pages}}
		Atom URL: <a href="{{.atom}}">{{.atom}}</a><br />
		<br />

		Selected pages:<br />
		{{- range .pages}}
			<input type="checkbox" name="p" value="{{.}}" id="p{{.}}" CHECKED /><label for="p{{.}}" title="{{.}}"> {{title .}}</label><br />
		{{- end}}
		<br />
	{{- end}}

	{{- if .available}}
		Add some pages we know about already:<br />
		{{- range .available}}
			<input type="checkbox" name="p" value="{{.}}" id="p{{.}}" /><label for="p{{.}}" title="{{.}}"> {{title .}}</label><br />
		{{- end}}
		<br />
	{{- end}}

	Or add other en.wikipedia.org pages (either the full URL or the part after <code>/wiki/</code>). One per line.<br />
	<textarea name="etc" cols="80" rows="4">
	</textarea><br />
	<br />

	<input type="submit" value="Update" /><br />
</form>

	<br />
	<hr />
	Version history of selected pages:<br />
	{{- range .versions}}
		{{- title .Page}}: {{.StableVersion}}<br />
	{{- end}}
{{- end}}
`))

	pageTempl = template.Must(extend(baseTempl).Parse(`
{{define "head"}}
	<link rel="alternate" type="application/atom+xml" title="Atom 1.0" href="{{.atom}}"/>
{{- end}}
{{define "page"}}
	{{.title}}<br />
	atom: <a href="{{.atom}}">{{.atom}}</a><br />
	<br />
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
