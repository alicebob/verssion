package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alicebob/w/w"
	"github.com/julienschmidt/httprouter"
)

var (
	dburl  = flag.String("db", "postgresql:///w", "postgres URL")
	listen = flag.String("listen", ":3141", "http listen")
)

func main() {
	flag.Parse()
	pages := flag.Args()
	if len(pages) != 0 {
		fmt.Fprintf(os.Stderr, "no args accepted\n")
		os.Exit(2)
	}

	db, err := w.NewPostgres(*dburl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg: %s\n", err)
		os.Exit(2)
	}

	up := newUpdate(db)

	r := httprouter.New()
	r.GET("/", indexHandler(db))
	r.GET("/v/:page/", pageHandler(db, up))
	r.GET("/v/:page/atom.xml", pageAtomHandler(db, up))
	fmt.Printf("listening on %s...\n", *listen)
	log.Fatal(http.ListenAndServe(*listen, r))
}

func indexHandler(db w.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		es, _ := db.Recent()
		runTmpl(w, indexTempl, map[string]interface{}{
			"title":   "hello world",
			"entries": es,
		})
	}
}

func pageHandler(db w.DB, up *update) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		page := p.ByName("page")
		if err := up.Update(page); err != nil {
			log.Printf("update %q: %s", page, err)
		}

		vs, _ := db.History(page)
		runTmpl(w, pageTempl, map[string]interface{}{
			"title":    page,
			"page":     page,
			"versions": vs,
		})
	}
}

func pageAtomHandler(db w.DB, up *update) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		page := p.ByName("page")
		if err := up.Update(page); err != nil {
			log.Printf("update %q: %s", page, err)
		}

		vs, err := db.History(page)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		if len(vs) == 0 {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		var (
			es   []Entry
			prev = "?"
		)
		for _, v := range vs {
			es = append(es, Entry{
				ID:      asuri(page, v.StableVersion),
				Title:   page + ": " + v.StableVersion,
				Updated: vs[len(vs)-1].T,
				Content: fmt.Sprintf("%s -> %s", prev, v.StableVersion),
			})
			prev = v.StableVersion
		}
		feed := Feed{
			XMLNS:   "http://www.w3.org/2005/Atom",
			ID:      asuri(page),
			Title:   page,
			Updated: vs[len(vs)-1].T,
			Author: Author{
				Name: "Wikipedia",
			},
			Entries: es,
		}
		w.Header().Set("Content-Type", "application/atom+xml")
		w.Write([]byte(xml.Header))
		e := xml.NewEncoder(w)
		e.Indent("", "\t")
		e.Encode(feed)
	}
}
