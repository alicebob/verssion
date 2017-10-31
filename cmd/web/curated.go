package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	libw "github.com/alicebob/w/w"
	"github.com/julienschmidt/httprouter"
)

func newCuratedHandler(db libw.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if r.FormValue("go") != "" {
			id, err := db.CreateCurated()
			if err != nil {
				log.Printf("create curated: %s", err)
				http.Error(w, http.StatusText(500), 500)
				return
			}
			w.Header().Set("Location", "./"+id+"/")
			w.WriteHeader(302)
			return
		}
		args := map[string]interface{}{
			"title": "curated list",
		}
		runTmpl(w, newCuratedTempl, args)
	}
}

func curatedHandler(db libw.DB, up *update, base string) httprouter.Handle {
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

		vs := map[string]libw.Entry{}
		for _, p := range cur.Pages {
			e, err := db.Current(p)
			if err == nil && e != nil {
				vs[p] = *e
			}
		}

		args := map[string]interface{}{
			"curated":      cur,
			"atom":         fmt.Sprintf("%s/curated/%s/atom.xml", base, cur.ID),
			"title":        cur.Title(),
			"pageversions": vs,
		}

		runTmpl(w, curatedTempl, args)
	}
}

func curatedEditHandler(db libw.DB, up *update, base string) httprouter.Handle {
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

		args := map[string]interface{}{
			"curated": cur,
		}

		if r.Method == "POST" {
			r.ParseForm()
			var (
				newPages = append(r.Form["p"], toPages(r.Form.Get("etc"))...)
				pages    []string
				errors   []string
			)
			for _, p := range newPages {
				if up != nil {
					if err := up.Update(p); err != nil {
						log.Printf("update %q: %s", p, err)
						errors = append(errors, fmt.Sprintf("%q: %s", p, err))
					} else {
						pages = append(pages, p)
					}
				} else {
					pages = append(pages, p)
				}
			}
			pages = unique(pages)
			cur.Pages = pages
			args["errors"] = errors
			if err := db.CuratedPages(cur.ID, pages); err != nil {
				log.Printf("curated pages: %s", err)
				http.Error(w, http.StatusText(500), 500)
				return
			}

			title := r.Form.Get("title")
			if err := db.CuratedTitle(cur.ID, title); err != nil {
				log.Printf("curated title: %s", err)
			}

			w.Header().Set("Location", "./")
			w.WriteHeader(302)
			return
		}
		args["defaulttitle"] = cur.DefaultTitle()
		args["customtitle"] = cur.CustomTitle

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

func curatedAtomHandler(db libw.DB, up *update, base string) httprouter.Handle {
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

		if up != nil {
			for _, p := range cur.Pages {
				if err := up.Update(p); err != nil {
					log.Printf("update %q: %s", p, err)
				}
			}
		}

		vs, err := db.History(cur.Pages...)
		if err != nil {
			log.Printf("history: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		feed := asFeed(base, "urn:uuid:"+cur.ID, cur.Title(), cur.LastUpdated, vs)
		feed.Links = []Link{
			{
				Href: fmt.Sprintf("%s/curated/%s/", base, cur.ID),
				Rel:  "alternate", // not strictly true...
				Type: "text/html",
			},
			{
				Href: fmt.Sprintf("%s/curated/%s/atom.xml", base, cur.ID),
				Rel:  "self",
				Type: "application/atom+xml",
			},
		}
		writeFeed(w, feed)

		if err := db.CuratedUsed(id); err != nil {
			log.Printf("curated used %q: %s", id, err)
		}
	}
}

var (
	newCuratedTempl = template.Must(extend(baseTempl).Parse(`
{{define "page"}}
	Curated list is a stored on the server. Whenever you change the list you don't have to update the RSS link, since that only has the list ID.<br />
	You can also share the link, and everyone can update the list.<br />
	<br />
	
	<form method="POST">
	<input type="submit" name="go" value="Start a list" />
	</form>
{{- end}}
`))

	curatedTempl = template.Must(extend(baseTempl).Parse(`
{{define "head"}}
	<link rel="alternate" type="application/atom+xml" title="Atom 1.0" href="{{.atom}}"/>
{{- end}}
{{define "page"}}
	<h2>{{.curated.Title}}</h2>
	Atom link: <a href="{{.atom}}">{{.atom}}</a><br />
	<br />
	{{- with .curated.Pages}}
		<table>
		<tr>
			<td>Page</td>
			<td>Stable Version</td>
			<td>Spider T</td>
		</tr>
		{{- range .}}
			{{$p := (index $.pageversions .)}}
			<tr>
			<td><a href="../../p/{{.}}/" title="{{.}}">{{title .}}</a></td>
			<td>{{$p.StableVersion}}</td>
			<td>{{$p.T}}</td>
			</tr>
		{{- end}}
		</table>
	{{- else}}
		No pages selected, yet.<br />
	{{- end}}
	<br />
	<a href="./edit.html">Edit this list</a><br />
	<br />
	<br />
{{end}}
`))

	curatedEditTempl = template.Must(extend(baseTempl).Parse(`
{{define "page"}}
	<h2>{{.curated.Title}}</h2>
	<br />
	<br />
	<form method="POST">
	Title: <input type="text" size="40" name="title" value="{{.customtitle}}" placeholder="{{.defaulttitle}}" /><br />
	Pages:<br />
	{{- with .curated.Pages}}
		{{- range .}}
			<input type="checkbox" name="p" value="{{.}}" id="p{{.}}" CHECKED /><label for="p{{.}}" title="{{.}}"> {{title .}}</label><br />
		{{- end}}
	{{- else}}
		No pages selected, yet.<br />
	{{- end}}
	<br />

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
{{end}}
`))
)
