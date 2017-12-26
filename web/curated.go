package web

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/alicebob/verssion/core"
)

func newCuratedHandler(base string, db core.DB, spider core.Spider) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		r.ParseForm()
		var (
			etc   = r.Form.Get("etc")
			pages = r.Form["p"]
			title = r.Form.Get("title")
		)
		pm := map[string]bool{}
		for _, p := range pages {
			pm[p] = true
		}
		args := map[string]interface{}{
			"title":        "New feed",
			"current":      "curated",
			"etc":          etc,
			"selected":     pm,
			"defaulttitle": "",
			"customtitle":  title,
		}
		if r.Method == "POST" {
			pages, errors := readPageArgs(db, spider, pages, etc)
			if len(pages) > 0 && len(errors) == 0 {
				id, err := db.CreateCurated()
				if err != nil {
					log.Printf("create curated: %s", err)
					http.Error(w, http.StatusText(500), 500)
					return
				}
				if err := db.CuratedSetPages(id, pages); err != nil {
					log.Printf("curated pages: %s", err)
				}
				if err := db.CuratedSetTitle(id, title); err != nil {
					log.Printf("curated title: %s", err)
				}

				w.Header().Set("Location", "./"+id+"/")
				w.WriteHeader(302)
				return
			}
			args["errors"] = errors
		}
		avail, err := db.Known()
		if err != nil {
			log.Printf("known: %s", err)
		}
		args["available"] = avail
		runTmpl(w, newCuratedTempl, args)
	}
}

func curatedHandler(base string, db core.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		cur, err := db.LoadCurated(id)
		if err != nil {
			log.Printf("load curated: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		if cur == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		vs, err := db.CurrentIn(cur.Pages...)
		if err != nil {
			log.Printf("current: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		args := map[string]interface{}{
			"title":        cur.Title(),
			"current":      "curated",
			"curated":      cur,
			"atom":         fmt.Sprintf("%s/curated/%s/atom.xml", base, id),
			"pageversions": vs,
		}

		c := &http.Cookie{
			Name:     "curated-" + id,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(30 * 24 * time.Hour),
		}
		w.Header().Add("Set-Cookie", c.String())

		runTmpl(w, curatedTempl, args)
	}
}

func curatedEditHandler(base string, db core.DB, spider core.Spider) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		cur, err := db.LoadCurated(id)
		if err != nil {
			log.Printf("load curated: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		if cur == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		r.ParseForm()
		var (
			etc    = r.Form.Get("etc")
			qPages = r.Form["p"]
		)
		selected := map[string]bool{}
		for _, p := range cur.Pages {
			selected[p] = true
		}
		args := map[string]interface{}{
			"title":        cur.Title(),
			"current":      "",
			"curated":      cur,
			"etc":          etc,
			"selected":     selected,
			"pages":        cur.Pages,
			"defaulttitle": cur.DefaultTitle(),
			"customtitle":  cur.CustomTitle,
		}
		if r.Method == "POST" {
			pages, errors := readPageArgs(db, spider, qPages, etc)
			title := r.Form.Get("title")
			args["customtitle"] = title
			if len(errors) == 0 {
				if err := db.CuratedSetPages(id, pages); err != nil {
					log.Printf("curated pages: %s", err)
					http.Error(w, http.StatusText(500), 500)
					return
				}

				if err := db.CuratedSetTitle(id, title); err != nil {
					log.Printf("curated title: %s", err)
				}

				w.Header().Set("Location", "./")
				w.WriteHeader(302)
				return
			}

			selected := map[string]bool{}
			for _, p := range qPages {
				selected[p] = true
			}
			args["selected"] = selected
			args["errors"] = errors
		}

		seen := map[string]struct{}{}
		for _, p := range cur.Pages {
			seen[p] = struct{}{}
		}
		var av []string
		if avail, err := db.Known(); err == nil {
			for _, p := range avail {
				if _, ok := seen[p]; !ok {
					av = append(av, p)
				}
			}
		}
		args["available"] = av
		runTmpl(w, curatedEditTempl, args)
	}
}

func curatedAtomHandler(base string, db core.DB, spider core.Spider) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		cur, err := db.LoadCurated(id)
		if err != nil {
			log.Printf("load curated: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		if cur == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		actualPages, errs := runUpdates(db, spider, cur.Pages)
		if len(errs) != 0 {
			log.Printf("curated atom runUpdates: %s", errs)
		}

		vs, err := db.History(actualPages...)
		if err != nil {
			log.Printf("history: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		feed := asFeed(base, "urn:uuid:"+id, cur.Title(), cur.LastUpdated, vs)
		feed.Links = []Link{
			{
				Href: fmt.Sprintf("%s/curated/%s/", base, id),
				Rel:  "alternate", // not strictly true...
				Type: "text/html",
			},
			{
				Href: fmt.Sprintf("%s/curated/%s/atom.xml", base, id),
				Rel:  "self",
				Type: "application/atom+xml",
			},
		}
		writeFeed(w, feed)

		if err := db.CuratedSetUsed(id); err != nil {
			log.Printf("curated used %q: %s", id, err)
		}
	}
}

var (
	newCuratedTempl = withPage(`
{{define "hero"}}
<h1>New Feed</h1>
<p>Create a new RSS/Atom feed combining multiple pages. You can always change the feed later.</p>
{{end}}

{{define "body"}}
{{template "errors" .errors}}
	
<form method="POST" class="form start-list-form">
{{template "pageselection" .}}
<input type="submit" name="go" class="btn" value="Start a list" />
</form>
{{- end}}
`)

	curatedTempl = withPage(`
{{define "head"}}
	<link rel="alternate" type="application/atom+xml" title="Atom 1.0" href="{{.atom}}"/>
{{- end}}

{{define "hero"}}
	<h2>{{.curated.Title}}</h2>
	<p>Atom link: <a href="{{.atom}}">{{.atom}}</a></p>
{{end}}

{{define "body"}}
{{- with .pageversions}}
	<table class="table-responsive table-2">
	<tr>
		<th class="hidden-xs">Page</th>
		<th class="hidden-xs">Stable version</th>
		<th class="hidden-xs">Spider timestamp</th>
	</tr>
	{{- range .}}
		<tr>
		<td><a href="/p/{{.Page}}/" title="{{.Page}}">{{title .Page}}</a></td>
		<td>{{version .StableVersion}}</td>
		<td class="hidden-xs">{{.T.Format "2006-01-02 15:04 UTC"}}</td>
		</tr>
	{{- end}}
	</table>
{{- else}}
<p>
	No pages selected, yet.<br />
</p>
{{- end}}

<p>
	<a href="./edit.html" class="btn">Edit this feed</a>
</p>
{{end}}
`)

	curatedEditTempl = withPage(`
{{define "hero"}}
<h1>{{.curated.Title}}</h1>
{{end}}

{{define "body"}}
{{template "errors" .errors}}
	
<form method="POST" class="form start-list-form">
{{template "pageselection" .}}
<input type="submit" class="btn" value="Update" /><br />
</form>
{{- end}}
`)
)

// read p and etc arguments
func readPageArgs(
	db core.DB,
	spider core.Spider,
	pages []string,
	etc string,
) ([]string, []error) {
	pages = append(pages, toPages(etc)...)

	finalPages, errors := runUpdates(db, spider, pages)

	return unique(finalPages), errors
}
