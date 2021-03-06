package core

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

type Table struct {
	Rows [][]string
}

// FindTables returns all top-level tables
func FindTables(doc *html.Node) ([]Table, error) {
	var ts []Table
	var f func(*html.Node) error
	f = func(n *html.Node) error {
		if n.Type == html.ElementNode && n.Data == "table" {
			t, err := tTable(n)
			if err != nil {
				return err
			}
			if t != nil {
				ts = append(ts, *t)
			}
			return nil
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := f(c); err != nil {
				return err
			}
		}
		return nil
	}
	return ts, f(doc)
}

type Link struct {
	Rel   string
	Href  string
	Title string
}

// FindHeadLink returns all <link> nodes
func FindHeadLink(doc *html.Node) ([]Link, error) {
	var nodes []Link
	var f func(*html.Node) error
	f = func(n *html.Node) error {
		if n.Type == html.ElementNode && n.Data == "link" {
			nodes = append(nodes, Link{
				Rel:   getAttr("rel", n.Attr),
				Href:  getAttr("href", n.Attr),
				Title: getAttr("title", n.Attr),
			})
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := f(c); err != nil {
				return err
			}
		}
		return nil
	}
	return nodes, f(doc)
}

func tTable(n *html.Node) (*Table, error) {
	tab := &Table{}
	var f func(*html.Node)
	f = func(n *html.Node) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "tr" {
				if r := tRow(c); r != nil {
					tab.Rows = append(tab.Rows, r)
				}
				continue
			}
			f(c)
		}
	}
	f(n)
	return tab, nil
}

func tRow(n *html.Node) []string {
	var row []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && (c.Data == "td" || c.Data == "th") {
				row = append(row, cleanSpace(tString(c)))
				continue
			}
			f(c)
		}
	}
	f(n)
	return row
}

func tString(n *html.Node) string {
	res := ""
node:
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch c.Type {
		case html.TextNode:
			res += stripCtl(c.Data)
		case html.ElementNode:
			switch c.Data {
			case "sup":
				continue node
			case "br":
				res += "\n"
			case "a":
				if ignoreTag(c.Data, c.Attr) {
					continue node
				}
				var (
					href = getAttr("href", c.Attr)
					name = tString(c)
				)
				if name == `[±]` {
					// special case for wikipedia 'edit' link
					continue
				}
				if strings.HasPrefix(href, "//") {
					href = "https:" + href
				}
				if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
					n := strings.NewReplacer(`[`, `\[`, `]`, `\]`).Replace(name)
					res += fmt.Sprintf("[%s](%s)", n, href)
				} else {
					res += name
				}
			default:
				if !ignoreTag(c.Data, c.Attr) {
					res += tString(c)
				}
				if c.Data == "tr" {
					res += "\n"
				}
			}
		default:
			res += tString(c)
		}
	}
	return res
}

// ignore certain HTML elemens. Specific to wikipedia
func ignoreTag(tag string, attr []html.Attribute) bool {
	for _, a := range attr {
		switch {
		case a.Key == "class" && strings.Contains(a.Val, "noprint"):
			return true
		case a.Key == "style" && strings.Contains(a.Val, "display:none"):
			return true
		case a.Key == "title" && strings.Contains(a.Val, "Wikidata"):
			return true
		}
	}
	return false
}

func stripCtl(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case unicode.IsSpace(r), unicode.IsControl(r):
			return ' '
		default:
			return r
		}
	}, s)
}

// Clean spaces in a multiline string. Keeps newlines.
func cleanSpace(s string) string {
	var res []string
	for _, l := range strings.Split(s, "\n") {
		if line := cleanLine(l); len(line) > 0 {
			res = append(res, line)
		}
	}
	return strings.Join(res, "\n")
}

func cleanLine(s string) string {
	var inSpace = false
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			if inSpace {
				return -1
			}
			inSpace = true
			return ' '
		}
		inSpace = false
		return r
	}, strings.TrimSpace(s))
}

func getAttr(attr string, attrs []html.Attribute) string {
	for _, a := range attrs {
		if a.Key == attr {
			return a.Val
		}
	}
	return ""
}
