package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	libw "github.com/alicebob/w/w"
	"github.com/julienschmidt/httprouter"
)

func pageHandler(db libw.DB, up *update) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		page := p.ByName("page")
		cur, err := up.Fetch(page, 10)
		if err != nil {
			if p, ok := err.(libw.ErrNotFound); ok {
				log.Printf("not found %q: %s", page, err)
				w.WriteHeader(404)
				runTmpl(w, pageNotFoundTempl, map[string]interface{}{
					"title":     libw.Title(p.Page),
					"wikipedia": libw.WikiURL(p.Page),
					"page":      p.Page,
				})
				return
			}
			log.Printf("update %q: %s", page, err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		if page != cur.Page {
			w.Header().Set("Location", fmt.Sprintf("%s/p/%s/", *baseURL, cur.Page))
			w.WriteHeader(302)
			return
		}

		vs, err := db.History(cur.Page)
		if err != nil {
			log.Printf("history: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		runTmpl(w, pageTempl, map[string]interface{}{
			"title":     libw.Title(cur.Page),
			"atom":      adhocURL([]string{cur.Page}),
			"wikipedia": libw.WikiURL(cur.Page),
			"current":   cur,
			"page":      cur.Page,
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

	pageNotFoundTempl = template.Must(extend(baseTempl).Parse(`
{{define "page"}}
	Page not found: {{.page}}<br />
	Maybe you can create it on <a href="{{.wikipedia}}">Wikipedia</a><br />
    <br />
{{- end}}
`))
)
