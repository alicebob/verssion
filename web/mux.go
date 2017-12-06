package web

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/alicebob/verssion/core"
)

func Mux(baseURL string, db core.DB, up Fetcher, static string) *httprouter.Router {
	r := httprouter.New()
	r.GET("/", indexHandler(baseURL, db))
	r.GET("/adhoc/atom.xml", adhocAtomHandler(baseURL, db, up))
	r.GET("/curated/", newCuratedHandler(baseURL, db, up))
	r.POST("/curated/", newCuratedHandler(baseURL, db, up))
	r.GET("/curated/:id/", curatedHandler(baseURL, db))
	r.GET("/curated/:id/edit.html", curatedEditHandler(baseURL, db, up))
	r.POST("/curated/:id/edit.html", curatedEditHandler(baseURL, db, up))
	r.GET("/curated/:id/atom.xml", curatedAtomHandler(baseURL, db, up))
	r.GET("/p/", allPagesHandler(baseURL, db))
	r.GET("/p/:page/", pageHandler(baseURL, db, up))
	if static != "" {
		r.ServeFiles("/s/*filepath", http.Dir(static))
	}
	return r
}
