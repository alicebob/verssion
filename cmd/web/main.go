package main

import (
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
	r.GET("/v/:page", pageHandler(db, up))
	fmt.Printf("listening on %s...\n", *listen)
	log.Fatal(http.ListenAndServe(*listen, r))
}

func indexHandler(db w.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		es, _ := db.Recent()
		runTmpl(w, indexTempl, map[string]interface{}{
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
			"page":     page,
			"versions": vs,
		})
	}
}
