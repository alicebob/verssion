package web

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/alicebob/verssion/core"
)

var (
	baseTempl = template.Must(
		template.New("base").
			Funcs(template.FuncMap{
				"title": core.Title,
				"version": func(s string) template.HTML {
					h := textMarkdown(s)
					h = template.HTMLEscapeString(h)
					return template.HTML(strings.Replace(h, "\n", "<br />", -1))
				},
				"link":      htmlMarkdown,
				"plainText": textMarkdown,
			}).Parse(`<!DOCTYPE html>
<html>
<head>
	<title>{{ .title }}</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="shortcut icon" href="{{.base}}/s/favicon.png" type="image/png" sizes="16x16 24x24 32x32 64x64">
	<link rel="apple-touch-icon" href="{{.base}}/s/favicon.png">
	<style type="text/css">
body {
	margin: 0 0 2em 0;
	font-family: "Helvetica Neue",Helvetica,Arial,sans-serif;
}
table {
	width: 100%;
	border-collapse: collapse;
}
th {
	text-align: left;
	font-weight: normal;
}
td, th {
    vertical-align: top;
	padding: 0;
	padding-bottom: 3px;
}
td:first-child {
	width: 250px;
}
textarea {
	box-sizing: border-box;
	width: 100%;
}
h2 {
	border-bottom: 1px solid #ddd;
}
a, a:visited {
	color: #357cb7;
}
a:hover {
	color: black;
}
nav {
	background-color: #35b7b1;
	padding: 0.5em 0;
}
nav ul {
	padding: 0;
	margin: 0;
}
nav li {
	list-style: none;
	display: inline;
	padding-right: 1em;
}
nav a {
	font-weight: bold;
}
nav a, nav a:visited {
	color: rgba(0, 0, 0, 0.5);
	text-decoration: none;
}
nav a.active, nav a:hover {
	color: rgba(0, 0, 0, 1);
}
.body, nav div {
	margin: 0 auto;
	max-width: 760px;
	padding: 0 0.5em;
}
.body p {
	text-align: justify;
}

@media only screen and (max-width: 700px) {
	table, thead, tbody, tr, th, td {
		display: block;
	}

	.optional {
		display: none;
	}
}
	</style>
	{{- block "head" .}}{{end}}
</head>
<body>
	<nav>
		<div>
		<ul>
        <li><a href="{{.base}}/"{{if eq .current "home"}} class="active"{{end}}>Home</a></li>
        <li><a href="{{.base}}/curated/"{{if eq .current "curated"}} class="active"{{end}}>New feed</a></li>
        <li><a href="{{.base}}/p/"{{if eq .current "allpages"}} class="active"{{end}}>All pages</a></li>
		</ul>
		</div>
	</nav>
	<div class="body">
        {{- block "page" .}}{{end}}
	</div>
</body>
</html>

{{define "errors"}}
    {{- with .}}
        Some problems:<br />
        {{- range .}}
            {{.}}<br />
        {{- end}}
        <br />
        <br />
    {{- end}}
{{end}}

{{define "pageselection"}}
	<script>
function moveChecked() {
	var avail = document.getElementById("available").childNodes;
	avail.forEach(function(elem) {
		if (elem.nodeType != Node.ELEMENT_NODE) {
			return
		};
		if (elem.firstElementChild.checked) {
			document.getElementById("selected").appendChild(elem);
		}
	});
}
function runFilter(v) {
    v = v.toLowerCase();
	var avail = document.getElementById("available").childNodes;
	avail.forEach(function(elem) {
		if (elem.nodeType != Node.ELEMENT_NODE) {
			return
		};
		if (! elem.dataset) {
			return;
		}
		var page = elem.dataset["page"].toLowerCase(),
			title = elem.dataset["title"].toLowerCase();
		if (title.length == 0) {
			return;
		};
		var visible = v.length == 0 || title.indexOf(v) >= 0 || page.indexOf(v) >= 0;
		elem.style.display = visible ? "block" : "none";
	});
}
	</script>
	<div id="selected">
    {{- if .pages}}
    Selected pages:<br />
    {{- range .pages}}
		<div>
        <input type="checkbox" name="p" value="{{.}}" id="p{{.}}"{{if (index $.selected .)}} CHECKED{{end}}/><label for="p{{.}}" title="{{.}}"> {{title .}}</label>
		</div>
    {{- end}}
    {{- end}}
	</div>
    <br />

    Add some pages:<br />
	Filter: <input type="text" oninput="moveChecked();runFilter(this.value)"><br />
	<div id="available">
    {{- range .available}}
		<div data-page="{{.}}" data-title="{{title .}}">
        <input type="checkbox" name="p" value="{{.}}" id="p{{.}}"{{if (index $.selected .)}} CHECKED{{end}}/><label for="p{{.}}" title="{{.}}"> {{title .}}</label>
		</div>
    {{- end}}
	</div>
    <br />

    Or add other en.wikipedia.org pages (either the full URL or the part after <code>/wiki/</code>). One per line.<br />
    <textarea name="etc" rows="4">{{.etc}}</textarea><br />
{{end}}
`))
)

func withBase(s string) *template.Template {
	return template.Must(template.Must(baseTempl.Clone()).Parse(s))
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
