package web

import (
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/alicebob/verssion/core"
)

func indexHandler(base string, db core.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		es, err := db.Current(12, core.SpiderT)
		if err != nil {
			log.Printf("current: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		curated, err := readCuratedCookies(r, db)
		if err != nil {
			log.Printf("readCuratedCookies: %s", err)
		}
		runTmpl(w, indexTempl, map[string]interface{}{
			"base":    base,
			"title":   "VeRSSion",
			"current": "home",
			"entries": es,
			"curated": curated,
		})
	}
}

func readCuratedCookies(r *http.Request, db core.DB) (map[string]*core.Curated, error) {
	var (
		curs    = map[string]*core.Curated{}
		lastErr error
	)
	for _, cookie := range r.Cookies() {
		t := strings.SplitN(cookie.Name, "-", 2)
		if len(t) != 2 || t[0] != "curated" {
			continue
		}
		id := t[1]
		c, err := db.LoadCurated(id)
		if err != nil {
			lastErr = err
		} else {
			if c != nil {
				curs[id] = c
			}
		}
	}
	return curs, lastErr
}

var (
	indexTempl = withBase(`
{{define "page"}}
<h2>What</h2>
<p>Verssion tracks stable version of software projects (e.g.: databases, editors, JS frameworks), and makes that available as an RSS (atom) feed. The main use-case is for dev-ops and developers who use a lot of open source software projects, and who like to keep an eye on releases. Without making that a fulltime job, and without signing up for dozens of e-mail lists. Turns out wikipedia is a great source for version information, so that's what we use.</p>
<p>You can create feeds for your own use, or share them with colleagues.</p>
<p><a href="https://github.com/alicebob/verssion/">Full source</a> for issues and suggestions.</p>

<h2>Feed</h2>
<p>Make a feed which combines multiple projects in a single feed:
<a href="./curated/">Create new feed!</a></p>

{{- if .curated}}
	<p>
	Your recent feeds:
	<ul>
	{{- range $id, $cur := .curated}}
		<li><a href="{{$.base}}/curated/{{$id}}/">{{$cur.Title}}</a></li>
	{{- end}}
	</ul>
	</p>
{{- end}}

<h2>Updates</h2>
	<table>
	{{- range .entries}}
		<tr>
			<td><a href="./p/{{.Page}}/">{{title .Page}}</a></td>
			<td>{{version .StableVersion}}</td>
		</tr>
	{{- end}}
		<tr>
			<td><a href="./p/?order=spider">...</a></td>
			<td></td>
		</tr>
	</table>

{{- end}}
`)
)
