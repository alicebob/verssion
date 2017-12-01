package web

import (
	"github.com/julienschmidt/httprouter"

	w "github.com/alicebob/verssion/w"
)

func Mux(baseURL string, db w.DB, up Fetcher) *httprouter.Router {
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
	return r
}
