package main

import (
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"

	libw "github.com/alicebob/verssion/w"
)

func adhocAtomHandler(db libw.DB, up *update, base string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		pages := r.URL.Query()["p"]
		sort.Strings(pages)
		actualPages, _ := runUpdates(up, pages)

		vs, err := db.History(actualPages...)
		if err != nil {
			log.Printf("history: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		feed := asFeed(
			base,
			asURN(strings.Join(actualPages, ",")),
			strings.Join(libw.Titles(actualPages), ", "),
			time.Time{},
			vs,
		)
		feed.Links = []Link{
			{
				Href: adhocURL(actualPages),
				Rel:  "self",
				Type: "application/atom+xml",
			},
		}
		writeFeed(w, feed)
	}
}
