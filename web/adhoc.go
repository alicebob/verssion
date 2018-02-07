package web

import (
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/alicebob/verssion/core"
)

func adhocAtomHandler(base string, db core.DB, spider core.Spider) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		pages := r.URL.Query()["p"]
		sort.Strings(pages)
		actualPages, errs := runUpdates(db, spider, pages)
		if len(errs) != 0 {
			log.Printf("adhoc atom runUpdates: %s", errs)
		}
		if len(actualPages) == 0 {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		vs, err := db.History(actualPages...)
		if err != nil {
			log.Printf("history: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		feed := asFeed(
			base,
			asURN(strings.Join(actualPages, ",")),
			strings.Join(core.Titles(actualPages), ", "),
			time.Time{},
			vs,
		)
		feed.Links = []Link{
			{
				Href: adhocURL(base, actualPages),
				Rel:  "self",
				Type: "application/atom+xml",
			},
		}
		writeFeed(w, feed)
	}
}
