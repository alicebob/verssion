package web

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/alicebob/verssion/core"
)

func allPages(base string, db core.DB, w http.ResponseWriter, r *http.Request) {
	all, err := db.CurrentAll()
	if err != nil {
		log.Printf("current all: %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	runTmpl(w, allPagesTempl, map[string]interface{}{
		"base":    base,
		"title":   "Pages overview",
		"current": "allpages",
		"pages":   all,
	})
}

// pageHandler deal with everything under `/p/`.
func pageHandler(base string, db core.DB, spider core.Spider) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		page := p.ByName("page")
		if page == "" || page == "/" {
			allPages(base, db, w, r)
			return
		}
		page = page[1:] // prefix /
		if !strings.HasSuffix(page, "/") {
			w.Header().Set("Location", base+"/p/"+page+"/")
			w.WriteHeader(301)
			return
		}
		page = page[0 : len(page)-1] // post /
		cur, err := StoreSpider(db, spider, page)
		if err != nil {
			if p, ok := err.(core.ErrNotFound); ok {
				log.Printf("not found: %s", err)
				w.WriteHeader(404)
				runTmpl(w, pageNotFoundTempl, map[string]interface{}{
					"base":      base,
					"title":     core.Title(p.Page),
					"current":   "",
					"wikipedia": WikiURL(p.Page),
					"page":      p.Page,
				})
				return
			}
			log.Printf("update %q: %s", page, err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		if page != cur.Page {
			w.Header().Set("Location", fmt.Sprintf("%s/p/%s/", base, cur.Page))
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
			"base":      base,
			"title":     core.Title(cur.Page),
			"current":   "",
			"atom":      adhocURL(base, []string{cur.Page}),
			"wikipedia": WikiURL(cur.Page),
			"page":      cur,
			"versions":  vs,
		})
	}
}

var (
	allPagesTempl = withBase(`
{{define "page"}}
	<h2>All known pages</h2>
	<table>
	{{- range .pages}}
		<tr>
			<td><a href="./{{.Page}}/" title="{{.Page}}">{{title .Page}}</a></td>
			<td>{{version .StableVersion}}</td>
		</tr>
	{{- end}}
	</table>
	<br />
	<a href="{{.base}}/curated/">Make a custom feed</a><br />
{{- end}}
`)

	pageTempl = withBase(`
{{define "head"}}
	<link rel="alternate" type="application/atom+xml" title="Atom 1.0" href="{{.atom}}"/>
{{- end}}
{{define "page"}}
	<h2>{{title .page.Page}}</h2>
	<table>
		<tr>
			<td>Wikipedia:</td>
			<td><a href="{{.wikipedia}}">{{.wikipedia}}</a></td>
		</tr>
		<tr>
			<td>Homepage:</td>
			<td>{{with .page.Homepage}}{{link .}}{{- end}}</td>
		</tr>
		<tr>
			<td>Current stable version:</td>
			<td>{{with .page.StableVersion}}{{version .}}{{- end}}</td>
		</tr>
	</table>
	<br />
	<br />
	<h2>Version history</h2>
	<table class="history">
	<tr>
		<th class="optional">Spider timestamp:</th>
		<th class="optional">Version:</th>
	</tr>
	{{- range .versions}}
		<tr>
			<td class="optional">{{.T.Format "2006-01-02 15:04 UTC"}}</td>
			<td>{{version .StableVersion}}</td>
		</tr>
	{{- end}}
	</table>
	<br />
	RSS link: <a href="{{.atom}}">Atom feed</a><br />
	<br />
	<small>
		Version numbers are retrieved from Wikipedia, and are licensed under Creative Commons.<br />
		If the current stable version is out of date, please edit <a href="{{.wikipedia}}">Wikipedia</a>.<br />
		Latest spider check: {{if not .page.T.IsZero}}{{.page.T.Format "2006-01-02 15:04 UTC"}}{{- end}}<br />
	</small>
{{- end}}
`)

	pageNotFoundTempl = withBase(`
{{define "page"}}
	Page not found: {{.page}}<br />
	Maybe you can create it on <a href="{{.wikipedia}}">Wikipedia</a><br />
	<br />
{{- end}}
`)
)
