package main

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
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

func asURN(s string) string {
	n := sha1.Sum([]byte(s))
	return fmt.Sprintf("urn:sha1:%x", n)
}
