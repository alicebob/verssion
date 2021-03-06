package web

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/alicebob/verssion/core"
)

func Mux(baseURL string, db core.DB, sp core.Spider, static string) *httprouter.Router {
	r := httprouter.New()
	r.GET("/", indexHandler(baseURL, db))
	r.GET("/adhoc/atom.xml", adhocAtomHandler(baseURL, db, sp))
	r.GET("/curated/", newCuratedHandler(baseURL, db, sp))
	r.POST("/curated/", newCuratedHandler(baseURL, db, sp))
	r.GET("/curated/:id/", curatedHandler(baseURL, db))
	r.GET("/curated/:id/edit.html", curatedEditHandler(baseURL, db, sp))
	r.POST("/curated/:id/edit.html", curatedEditHandler(baseURL, db, sp))
	r.GET("/curated/:id/atom.xml", curatedAtomHandler(baseURL, db, sp))
	r.GET("/p/*page", pageHandler(baseURL, db, sp))
	r.GET("/robots.txt", robotsHandler)
	if static != "" {
		r.ServeFiles("/s/*filepath", http.Dir(static))
	}
	return r
}
