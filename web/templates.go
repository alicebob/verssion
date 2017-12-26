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

    <link rel="apple-touch-icon" sizes="180x180" href="/s/favicons/apple-touch-icon.png">
    <link rel="icon" type="image/png" sizes="32x32" href="/s/favicons/favicon-32x32.png">
    <link rel="icon" type="image/png" sizes="16x16" href="/s/favicons/favicon-16x16.png">
    <link rel="manifest" href="/s/favicons/manifest.json">
    <link rel="mask-icon" href="/s/favicons/safari-pinned-tab.svg" color="#5bbad5">

    <link rel="stylesheet" href="/s/fonts/font-awesome.css">
    <link rel="stylesheet" href="/s/verssion.css">
    <link rel="stylesheet" href="/s/grid.css">

	{{- block "head" .}}{{end}}
</head>
<body>

<nav role="navigation">
    <ul>
    <li class="logo"><a href="/" title="Verssion"></a></li>
    <li{{if eq .current "home"}} class="current-menu-item"{{end}}><a href="/">Home</a></li>
    <li{{if eq .current "curated"}} class="current-menu-item"{{end}}><a href="/curated/">New Feed</a></li>
    <li{{if eq .current "allpages"}} class="current-menu-item"{{end}}><a href="/p/">All Pages</a></li>
    </ul>
</nav>

<div class="herobg">
</div>
{{- block "page" .}}{{end}}

<footer>
	<p><a href="https://github.com/alicebob/verssion"><strong>VERSSION</strong> <span class="fab fa-github" /></a> :: Version data <span class="fab fa-creative-commons"></span> Wikipedia</p>
</footer>

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
	
	<dl>

	<dt>Title</dt>
	<dd>
		<input type="text" name="title" value="{{.customtitle}}" placeholder="{{.defaulttitle}}" />
	</dd>

	<dt>Pages</dt>
	<dd>
		<div class="checks" id="selected">
		{{- range .pages}}
			<div>
				<input type="checkbox" name="p" value="{{.}}" id="p{{.}}"{{if (index $.selected .)}} CHECKED{{end}}/><label for="p{{.}}" title="{{.}}"> {{title .}}</label>
			</div>
		{{- end}}
		</div>
		<label class="input">filter</label>
		<input type="text" oninput="moveChecked();runFilter(this.value)">
		<div class="checks" id="available">
		{{- range .available}}
			<div data-title="{{title .}}">
				<input type="checkbox" name="p" value="{{.}}" id="p{{.}}" {{if (index $.selected .)}} CHECKED{{end}}/><label for="p{{.}}" title="{{.}}">{{title .}}</label>
			</div>
		{{- end}}
		</div>
	</dd>
		

	<dt>Add more pages</dt>
	<dd>
		<p>Or add other en.wikipedia.org pages (either the full URL or the part after <code>/wiki/</code>). One per line.</p>
		<textarea name="etc" rows="4">{{.etc}}</textarea>
	</dd>
	</dl>
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
