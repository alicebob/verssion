package main

import (
	"encoding/xml"
	"log"
	"net/http"
)

func writeFeed(w http.ResponseWriter, feed Feed) {
	w.Header().Set("Content-Type", "application/atom+xml")
	w.Write([]byte(xml.Header))
	e := xml.NewEncoder(w)
	e.Indent("", "\t")
	if err := e.Encode(feed); err != nil {
		log.Printf("xml encode: %s", err)
	}
}
