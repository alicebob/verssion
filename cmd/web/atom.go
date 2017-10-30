package main

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"time"

	libw "github.com/alicebob/w/w"
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

func asURN(s string) string {
	n := sha1.Sum([]byte(s))
	return fmt.Sprintf("urn:sha1:%x", n)
}

func asFeed(id, title string, update time.Time, vs []libw.Entry) Feed {
	var es []Entry
	for _, v := range vs {
		es = append(es, Entry{
			ID:      asURN(v.Page + "-" + v.StableVersion),
			Title:   libw.Title(v.Page) + ": " + v.StableVersion,
			Updated: v.T,
			Content: v.StableVersion, // TODO: prev version?
		})
		if v.T.After(update) {
			update = v.T
		}
	}
	return Feed{
		XMLNS:   "http://www.w3.org/2005/Atom",
		ID:      id,
		Title:   title,
		Updated: update,
		Author: Author{
			Name: "Wikipedia",
		},
		Entries: es,
	}
}
