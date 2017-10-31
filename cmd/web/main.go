package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"

	libw "github.com/alicebob/w/w"
	"github.com/julienschmidt/httprouter"
)

var (
	baseURL = flag.String("base", "http://localhost:3141", "base URL")
	dbURL   = flag.String("db", "postgresql:///w", "postgres URL")
	listen  = flag.String("listen", ":3141", "http listen")
	updates = flag.Bool("update", true, "update pages")
)

func main() {
	flag.Parse()
	pages := flag.Args()
	if len(pages) != 0 {
		fmt.Fprintf(os.Stderr, "no args accepted\n")
		os.Exit(2)
	}

	db, err := libw.NewPostgres(*dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pg: %s\n", err)
		os.Exit(2)
	}

	up := newUpdate(db)
	if !*updates {
		up = nil
	}

	r := httprouter.New()
	r.GET("/", indexHandler(db))
	r.GET("/adhoc/", adhocHandler(db, up))
	r.GET("/adhoc/atom.xml", adhocAtomHandler(db, up, *baseURL))
	r.GET("/curated/", newCuratedHandler(db))
	r.POST("/curated/", newCuratedHandler(db))
	r.GET("/curated/:id/", curatedHandler(db, up, *baseURL))
	r.POST("/curated/:id/", curatedHandler(db, up, *baseURL))
	r.GET("/curated/:id/atom.xml", curatedAtomHandler(db, up, *baseURL))
	r.GET("/p/:page/", pageHandler(db, up))
	fmt.Printf("listening on %s...\n", *listen)
	log.Fatal(http.ListenAndServe(*listen, r))
}

func indexHandler(db libw.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		es, _ := db.Recent()
		runTmpl(w, indexTempl, map[string]interface{}{
			"title":   "hello world",
			"entries": es,
		})
	}
}

func adhocURL(pages []string) string {
	u, err := url.Parse(*baseURL)
	if err != nil {
		panic(err)
	}
	u.Path += "/adhoc/atom.xml"
	u.RawQuery = url.Values{
		"p": pages,
	}.Encode()
	return u.String()
}

var matchpage = regexp.MustCompile(`^(?:(?i:https?://en.wikipedia.org)/wiki/)?(\S+)$`)

// from textarea to pages
func toPages(q string) []string {
	var ps []string
	for _, l := range strings.Split(q, "\n") {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		if m := matchpage.FindStringSubmatch(l); m != nil {
			ps = append(ps, m[1])
		}
	}
	return ps
}

func unique(ps []string) []string {
	m := map[string]struct{}{}
	for _, p := range ps {
		m[p] = struct{}{}
	}
	res := make([]string, 0, len(m))
	for p := range m {
		res = append(res, p)
	}
	sort.Strings(res)
	return res
}
