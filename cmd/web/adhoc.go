package main

import (
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	libw "github.com/alicebob/w/w"
	"github.com/julienschmidt/httprouter"
)

func adhocHandler(db libw.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		q := r.URL.Query()
		pages := q["p"]
		etcP, _ := toPages(q.Get("etc"))
		pages = append(pages, etcP...)
		pages = unique(pages)
		var (
			errors []string
			seen   = map[string]struct{}{}
		)
		for _, p := range pages {
			seen[p] = struct{}{}
		}

		args := map[string]interface{}{
			"errors": errors,
			"pages":  pages,
			"title":  "adhoc atom builder",
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

		if len(pages) > 0 {
			vs, err := db.History(pages...)
			if err != nil {
				log.Printf("history: %s", err)
				http.Error(w, http.StatusText(500), 500)
				return
			}
			args["title"] = strings.Join(libw.Titles(pages), ", ")
			args["versions"] = vs
			args["atom"] = adhocURL(pages)
		}
		runTmpl(w, adhocTempl, args)
	}
}

func adhocAtomHandler(db libw.DB, up *update, base string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		pages := r.URL.Query()["p"]
		sort.Strings(pages)
		if up != nil {
			for _, p := range pages {
				if err := up.Update(p); err != nil {
					log.Printf("update %q: %s", p, err)
				}
			}
		}

		vs, err := db.History(pages...)
		if err != nil {
			log.Printf("history: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		feed := asFeed(
			base,
			asURN(strings.Join(pages, ",")),
			strings.Join(libw.Titles(pages), ", "),
			time.Time{},
			vs,
		)
		feed.Links = []Link{
			{
				Href: adhocURL(pages),
				Rel:  "self",
				Type: "application/atom+xml",
			},
		}
		writeFeed(w, feed)
	}
}

var (
	adhocTempl = template.Must(extend(baseTempl).Parse(`
{{define "head"}}
	<link rel="alternate" type="application/atom+xml" title="Atom 1.0" href="{{.atom}}"/>
{{- end}}
{{define "page"}}

{{with .errors}}
	&ltblink>Errors&lt/blink><br />
	{{- range .}}
		{{.}}<br />
	{{- end}}
	<br />
{{end}}

<form method="GET">
	Create an ad hoc Atom URL to track versions.<br />
	<br />
	<br />
	{{- if .pages}}
		Atom URL: <a href="{{.atom}}">{{.atom}}</a><br />
		<br />

		Selected pages:<br />
		{{- range .pages}}
			<input type="checkbox" name="p" value="{{.}}" id="p{{.}}" CHECKED /><label for="p{{.}}" title="{{.}}"> {{title .}}</label><br />
		{{- end}}
		<br />
	{{- end}}

	{{- if .available}}
		Add some pages we know about already:<br />
		{{- range .available}}
			<input type="checkbox" name="p" value="{{.}}" id="p{{.}}" /><label for="p{{.}}" title="{{.}}"> {{title .}}</label><br />
		{{- end}}
		<br />
	{{- end}}

	Or add other en.wikipedia.org pages (either the full URL or the part after <code>/wiki/</code>). One per line.<br />
	<textarea name="etc" cols="80" rows="4">
	</textarea><br />
	<br />

	<input type="submit" value="Update" /><br />
</form>

	<br />
	<hr />
	Version history of selected pages:<br />
	{{- range .versions}}
		{{- title .Page}}: {{.StableVersion}}<br />
	{{- end}}
{{- end}}
`))
)
