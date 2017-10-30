package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

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
	r.GET("/adhoc/atom.xml", adhocAtomHandler(db, up))
	r.GET("/v/:page/", pageHandler(db, up))
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

func adhocHandler(db libw.DB, up *update) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		q := r.URL.Query()
		pages := q["p"]
		pages = append(pages, toPages(q.Get("etc"))...)
		pages, seen := unique(pages)
		var errors []string

		for _, p := range pages {
			if up != nil {
				if err := up.Update(p); err != nil {
					log.Printf("update %q: %s", p, err)
					errors = append(errors, fmt.Sprintf("%q: %s", p, err))
				}
			}
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

func adhocAtomHandler(db libw.DB, up *update) httprouter.Handle {
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

		var es []Entry
		for _, v := range vs {
			es = append(es, Entry{
				ID:      asURN(v.Page + "-" + v.StableVersion),
				Title:   libw.Title(v.Page) + ": " + v.StableVersion,
				Updated: v.T,
				Content: v.StableVersion, // TODO: prev version?
			})
		}
		var update time.Time
		if len(vs) > 0 {
			update = vs[len(vs)-1].T
		}
		writeFeed(w,
			asURN(strings.Join(pages, ",")),
			strings.Join(libw.Titles(pages), ", "),
			update,
			es,
		)
	}
}

func pageHandler(db libw.DB, up *update) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		page := p.ByName("page")
		if up != nil {
			if err := up.Update(page); err != nil {
				log.Printf("update %q: %s", page, err)
			}
		}

		vs, err := db.History(page)
		if err != nil {
			log.Printf("history: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		runTmpl(w, pageTempl, map[string]interface{}{
			"title":    libw.Title(page),
			"atom":     adhocURL([]string{page}),
			"page":     page,
			"versions": vs,
		})
	}
}

func writeFeed(w http.ResponseWriter, id, title string, update time.Time, es []Entry) {
	feed := Feed{
		XMLNS:   "http://www.w3.org/2005/Atom",
		ID:      id,
		Title:   title,
		Updated: update,
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

func unique(ps []string) ([]string, map[string]struct{}) {
	m := map[string]struct{}{}
	for _, p := range ps {
		m[p] = struct{}{}
	}
	res := make([]string, 0, len(m))
	for p := range m {
		res = append(res, p)
	}
	sort.Strings(res)
	return res, m
}
