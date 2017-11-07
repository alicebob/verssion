package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	libw "github.com/alicebob/verssion/w"
)

func indexHandler(db libw.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		es, err := db.Recent(12)
		if err != nil {
			log.Printf("current all: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		runTmpl(w, indexTempl, map[string]interface{}{
			"title":   "hello world",
			"entries": es,
		})
	}
}

var (
	indexTempl = template.Must(extend(baseTempl).Parse(`
{{define "page"}}
<div style="width: 40em"> 
Verssion(*) tracks stable version of software projects (e.g.: databases, editors, JS frameworks), and makes that available as an RSS (atom) feed. The main use-case is for dev-ops and developers who use a lot of open source software projects, and who like to keep an eye on releases. Without making that a fulltime job, and without signing up for dozen of email lists. Turns out wikipedia is a great source for version information, so that's what we use.<br />
You can create feeds for your own use, or share them with collegues.<br />
<br />
*) working title<br />
</div>
<br />

<a href="./curated/">New feed!</a><br />
<a href="./adhoc/">ad hoc feed</a> (kinda outdated)<br />
<br />

Some recent updates:<br />
	<table>
	{{- range .entries}}
		<tr>
			<td><a href="./p/{{.Page}}/">{{title .Page}}</a></td>
			<td>{{.StableVersion}}</td>
		</tr>
	{{- end}}
	</table>
{{- end}}
`))
)
