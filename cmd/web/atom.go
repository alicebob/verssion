package main

import (
	"encoding/xml"
	"net/url"
	"time"
)

type Feed struct {
	XMLName xml.Name  `xml:"feed"`
	XMLNS   string    `xml:"xmlns,attr"`
	ID      string    `xml:"id"`
	Title   string    `xml:"title"`
	Updated time.Time `xml:"updated"`
	Author  Author    `xml:"author"`
	Entries []Entry   `xml:"entry"`
}

type Author struct {
	Name string `xml:"name"`
}

type Entry struct {
	ID      string    `xml:"id"`
	Title   string    `xml:"title"`
	Updated time.Time `xml:"updated"`
	Content string    `xml:"content"`
}

func asuri(elems ...string) string {
	p := ""
	for _, s := range elems {
		p += "/" + url.PathEscape(s)
	}
	return "http://lijzij.de/w" + p
}
