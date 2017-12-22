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
					h := core.TextMarkdown(s)
					h = template.HTMLEscapeString(h)
					return template.HTML(strings.Replace(h, "\n", "<br />", -1))
				},
				"link":      core.HtmlMarkdown,
				"plainText": core.TextMarkdown,
			}).Parse(`<!DOCTYPE html>
<html>
<head>
	<title>{{ .title }}</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">

    <link rel="apple-touch-icon" sizes="180x180" href="{{.base}}/s/favicons/apple-touch-icon.png">
    <link rel="icon" type="image/png" sizes="32x32" href="{{.base}}/s/favicons/favicon-32x32.png">
    <link rel="icon" type="image/png" sizes="16x16" href="{{.base}}/s/favicons/favicon-16x16.png">
    <link rel="manifest" href="{{.base}}/s/favicons/manifest.json">
    <link rel="mask-icon" href="{{.base}}/s/favicons/safari-pinned-tab.svg" color="#5bbad5">

    <link rel="stylesheet" href="{{.base}}/s/bootstrap.css">
    <link rel="stylesheet" href="{{.base}}/s/pure-menu.css">
    <link rel="stylesheet" href="{{.base}}/s/pure-responsive-menu.css">
    <link rel="stylesheet" href="{{.base}}/s/fonts/font-awesome.css">
    <link rel="stylesheet" href="{{.base}}/s/verssion.css">
    <link rel="stylesheet" href="{{.base}}/s/grid.css">

	{{- block "head" .}}{{end}}
</head>
<body>
<div class="navtobe navbar-fixed-top" role="navigation" id="header">
<div class="container">
<div class="custom-wrapper pure-g" id="menu-container">
<div class="wrapper">
	<div id="brand-image-field">
		<div class="pure-menu">
			<a class="navbar-brand" title="Verssion" href="{{.base}}/">Verssion</a>
			<a style="z-index:8866;" href="#" class="custom-toggle" id="toggle">
				<s class="bar"></s>
				<s class="bar"></s>
				<s class="bar"></s>
			</a>
		</div>
	</div>
	<div id="menu-field">
		<nav id="menu" role="navigation" class="nav pure-menu pure-menu-horizontal custom-can-transform">
			<ul class="pure-menu-list" id="navigation-bar">
				<li class="menu-item{{if eq .current "home"}} current-menu-item{{end}}"><a href="{{.base}}/">Home</a></li>
				<li class="menu-item{{if eq .current "curated"}} current-menu-item{{end}}"><a href="{{.base}}/curated/">New Feed</a></li>
				<li class="menu-item{{if eq .current "allpages"}} current-menu-item{{end}}"><a href="{{.base}}/p/">All Pages</a></li>
			</ul>
		</nav>
	</div>
</div>
</div>
</div>
</div>

<div class="herobg">
</div>
{{- block "page" .}}{{end}}

<footer>
	<p>&copy; Copyright 2017 <i><strong>VERSSION</strong></i></p>
</footer>

<script src="{{.base}}/s/jquery.min.js"></script>
<script src="{{.base}}/s/scripts.js"></script>
</body>
</html>

{{define "errors"}}
	{{- with .}}
		<p>
		Some problems:
		<ul>
		{{- range .}}
			<li>{{.}}</li>
		{{- end}}
		</ul>
		</p>
	{{- end}}
{{end}}

{{define "pageselection"}}
	<script>
function moveChecked() {
	var avail = document.getElementById("available").childNodes;
	avail.forEach(function(elem) {
		if (elem.nodeType != Node.ELEMENT_NODE) {
			return;
		}
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
			return;
		}
		if (! elem.dataset) {
			return;
		}
		var title = elem.dataset["title"].toLowerCase();
		if (title.length == 0) {
			return;
		}
		var visible = v.length == 0 || title.indexOf(v) >= 0;
		elem.style.display = visible ? "block" : "none";
	});
}
	</script>

	{{- if .pages}}
	<h4>Selected pages</h4>
	{{- end}}
	<div class="row">
		<div class="col-sm-12">
			<div class="row checks" id="selected">
			{{- range .pages}}
				<div class="col-sm-12 check-row">
				<input type="checkbox" name="p" value="{{.}}" id="p{{.}}"{{if (index $.selected .)}} CHECKED{{end}}/><label for="p{{.}}" title="{{.}}"> {{title .}}</label>
				</div>
			{{- end}}
			</div>
		</div>
	</div>

	<h4>Known pages</h4>
	<div class="row">
		<div class="col-sm-12">
			<div class="row filter">
				<div class="col-sm-2"><label for="input-1">filter</label></div>
				<div class="col-sm-7">
				<input type="text" oninput="moveChecked();runFilter(this.value)">
				</div>
			</div>
			<div class="row checks" id="available">
			{{- range .available}}
				<div class="col-sm-12 check-row" data-title="{{title .}}">
					<input type="checkbox" name="p" value="{{.}}" id="p{{.}}" {{if (index $.selected .)}} CHECKED{{end}}/><label for="p{{.}}" title="{{.}}">{{title .}}</label>
				</div>
			{{- end}}
			</div>
		</div>
	</div>
		

	<div class="row">
		<div class="col-sm-12">
			<h4>Add more pages</h4>
			<p>Or add other en.wikipedia.org pages (either the full URL or the part after <code>/wiki/</code>). One per line.</p>
			<textarea name="etc" rows="4">{{.etc}}</textarea>
		</div>
	</div>
{{end}}
`))

	basePageTempl = `
{{define "page"}}
<div class="hero page">
{{- block "hero" .}}{{end}}
</div>

<div class="pagebg"></div>
<div class="content page">
{{- block "body" .}}{{end}}
</div>
{{end}}
`
)

func withBase(s string) *template.Template {
	return template.Must(template.Must(baseTempl.Clone()).Parse(s))
}

func withPage(s string) *template.Template {
	return template.Must(withBase(basePageTempl).Parse(s))
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
