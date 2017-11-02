package main

import (
	"html/template"
	"log"
	"net/http"

	libw "github.com/alicebob/w/w"
	"github.com/julienschmidt/httprouter"
)

func pageHandler(db libw.DB, up *update) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		page := p.ByName("page")
		if up != nil {
			if err := up.Update(page); err != nil {
				log.Printf("update %q: %s", page, err)
			}
		}

		curs, err := db.Current(page)
		if err != nil {
			log.Printf("current: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		cur := libw.Page{}
		for _, c := range curs {
			if c.Page == page {
				cur = c
				break
			}
		}

		vs, err := db.History(page)
		if err != nil {
			log.Printf("history: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		runTmpl(w, pageTempl, map[string]interface{}{
			"title":     libw.Title(page),
			"atom":      adhocURL([]string{page}),
			"wikipedia": libw.WikiURL(page),
			"current":   cur,
			"page":      page,
			"versions":  vs,
		})
	}
}

var (
	pageTempl = template.Must(extend(baseTempl).Parse(`
{{define "head"}}
	<link rel="alternate" type="application/atom+xml" title="Atom 1.0" href="{{.atom}}"/>
{{- end}}
{{define "page"}}
	{{.title}}<br />
	Atom feed: <a href="{{.atom}}">{{.atom}}</a><br />
	Wikipedia: <a href="{{.wikipedia}}">{{.wikipedia}}</a><br />
	Homepage: {{with .current.Homepage}}<a href="https://{{.}}">{{.}}</a>{{- end}}<br />
	Latest version: {{with .current.StableVersion}}{{.}}{{- end}}<br />
    <br />
    <br />
	History:
	<table>
	<tr>
		<th>Version Wiki text</th>
		<th>Spider T</th>
	</tr>
	{{- range .versions}}
		<tr>
			<td>{{- .StableVersion}}</td>
			<td>{{.T}})</td>
		</tr>
	{{- end}}
{{- end}}
`))
)
