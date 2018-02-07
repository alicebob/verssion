package web

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/alicebob/verssion/core"
)

type Feed struct {
	XMLName xml.Name  `xml:"feed"`
	XMLNS   string    `xml:"xmlns,attr"`
	ID      string    `xml:"id"`
	Title   string    `xml:"title"`
	Links   []Link    `xml:"link"`
	Updated time.Time `xml:"updated"`
	Author  Author    `xml:"author"`
	Entries []Entry   `xml:"entry"`
}

type Link struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
}

type Author struct {
	Name string `xml:"name"`
}

type Entry struct {
	ID      string    `xml:"id"`
	Title   string    `xml:"title"`
	Updated time.Time `xml:"updated"`
	Content string    `xml:"content"`
	Links   []Link    `xml:"link"`
}

func asURN(s string) string {
	n := sha1.Sum([]byte(s))
	return fmt.Sprintf("urn:sha1:%x", n)
}

func asFeed(base, id, title string, update time.Time, vs []core.Page) Feed {
	var es []Entry
	for _, v := range vs {
		sv := core.TextMarkdown(v.StableVersion)
		es = append(es, Entry{
			ID:      asURN(v.Page + "-" + v.StableVersion),
			Title:   core.Title(v.Page) + ": " + sv,
			Updated: v.T,
			Content: sv,
			Links: []Link{
				{
					Href: fmt.Sprintf("%s/p/%s/", base, v.Page),
					Rel:  "alternate", // not strictly true...
					Type: "text/html",
				},
			},
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
