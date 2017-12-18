package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/alicebob/verssion/core"
)

func allPages(base string, db core.DB, w http.ResponseWriter, r *http.Request) {
	if url := r.FormValue("page"); url != "" {
		// "new page" form
		if page := core.WikiBasePage(url); page != "" {
			w.Header().Set("Location", base+"/p/"+page+"/")
			w.WriteHeader(301)
			return
		}
	}
	qorder := r.FormValue("order")
	order := core.Alphabet
	if qorder == "spider" {
		order = core.SpiderT
	}
	all, err := db.Current(0, order)
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
		"order":   qorder,
	})
}

// PageJSON is returned for ?format=json
type PageJSON struct {
	Page          string `json:"page"`
	Title         string `json:"title"`
	StableVersion string `json:"stable_version"`
	Homepage      string `json:"homepage"`
	WikipediaURL  string `json:"wikipedia_url"`
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
		page = page[:len(page)-1] // post /
		cur, err := StoreSpider(db, spider, page)
		if err != nil {
			switch e := err.(type) {
			case core.ErrNotFound:
				log.Printf("not found: %s", err)
				w.WriteHeader(404)
				switch r.FormValue("format") {
				case "json":
				default:
					runTmpl(w, pageNotFoundTempl, map[string]interface{}{
						"base":      base,
						"title":     core.Title(e.Page),
						"current":   "",
						"wikipedia": WikiURL(e.Page),
						"page":      e.Page,
					})
				}
			case core.ErrNoVersion:
				switch r.FormValue("format") {
				case "json":
					w.WriteHeader(404)
				default:
					runTmpl(w, noVersionTempl, map[string]interface{}{
						"base":    base,
						"title":   core.Title(e.Page),
						"current": "",
						"page":    e.Page,
					})
				}
			default:
				log.Printf("update %q: %s", page, err)
				http.Error(w, http.StatusText(500), 500)
			}
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

		switch r.FormValue("format") {
		case "json":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(PageJSON{
				Page:          cur.Page,
				Title:         core.Title(cur.Page),
				WikipediaURL:  WikiURL(cur.Page),
				StableVersion: cur.StableVersion,
				Homepage:      cur.Homepage,
			})
		case "", "html":
			runTmpl(w, pageTempl, map[string]interface{}{
				"base":      base,
				"title":     core.Title(cur.Page),
				"current":   "",
				"atom":      adhocURL(base, []string{cur.Page}),
				"json":      fmt.Sprintf("%s/p/%s/?format=json", base, cur.Page),
				"wikipedia": WikiURL(cur.Page),
				"page":      cur,
				"versions":  vs,
			})
		default:
			http.Error(w, http.StatusText(400), 400)
		}
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
	<p>
	Order by: {{if eq .order "spider"}}
		Update - <a href="./">Alphabetical</a>
	{{else}}
		<a href="./?order=spider">Update</a> - Alphabetical
	{{end}}
	</p>
	<h4>New page</h4>
	<p>
		You're missing something? Add a wikipedia page. Give either the full URL or the title of the page.<br />
		<form method="GET">
		<input type="text" name="page" size="60" />
		<input type="submit" value="Go" />
		</form>
	</p>
{{- end}}
`)

	pageTempl = withBase(`
{{define "head"}}
	<link rel="alternate" type="application/atom+xml" title="Atom 1.0" href="{{.atom}}"/>
	<link rel="alternate" type="application/json" title="JSON" href="{{.json}}"/>
{{- end}}
{{define "page"}}
	<h2>{{title .page.Page}}</h2>
	<table>
		<tr>
			<td>Wikipedia</td>
			<td><a href="{{.wikipedia}}">{{.wikipedia}}</a></td>
		</tr>
		<tr>
			<td>Homepage</td>
			<td>{{with .page.Homepage}}{{link .}}{{- end}}</td>
		</tr>
		<tr>
			<td>Current stable version</td>
			<td>{{with .page.StableVersion}}{{version .}}{{- end}}</td>
		</tr>
	</table>
	<br />
	<br />
	<h2>Version history</h2>
	<table class="history">
	<tr>
		<th class="optional">Spider timestamp</th>
		<th class="optional">Version</th>
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
	<h2>{{.page}}</h2>
	Page not found: {{.page}}<br />
	Maybe you can create it on <a href="{{.wikipedia}}">Wikipedia</a><br />
	<br />
{{- end}}
`)

	noVersionTempl = withBase(`
{{define "page"}}
	<h2>{{.page}}</h2>
	No version found for <b>{{.page}}</b><br />
	Maybe this is a disambiguation page? <a href="https://en.wikipedia.org/w/index.php?search={{.page}}">Search</a> on Wikipedia.<br />
	<br />
{{- end}}
`)
)
