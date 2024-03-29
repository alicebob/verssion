package core

import (
	"strings"

	"golang.org/x/net/html"
)

func CanonicalPage(doc *html.Node) (string, error) {
	links, err := FindHeadLink(doc)
	if err != nil {
		return "", err
	}
	for _, link := range links {
		if link.Rel == "canonical" {
			return WikiBasePage(link.Href), nil
		}
	}
	return "", nil
}

func StableVersion(n *html.Node) (string, string) {
	var stable, homepage string

	ts, err := FindTables(n)
	if err != nil {
		return "", ""
	}
	for _, t := range ts {
		for i, r := range t.Rows {
			if len(r) == 0 {
				continue
			}
			v := ""
			if len(r) > 1 {
				v = r[1]
			}
			switch k := r[0]; TextMarkdown(k) {
			case "Stable release", "Latest release", "Last release":
				// vim has multiple projects on the page
				if stable == "" && v != "" {
					stable = v
				}
			case "Stable release(s)":
				// Firefox, has a table with versions. The version is in the
				// next row.
				if len(t.Rows) > i {
					if nextRow := t.Rows[i+1]; len(nextRow) > 0 {
						stable = nextRow[0]
					}
				}
			case "Official website", "Website":
				if homepage == "" && v != "" {
					homepage = v
				}
			}
		}
	}
	return stable, homepage
}

// title version of a wikipage path
func Title(page string) string {
	s := strings.Replace(page, "_", " ", -1)
	// remove some common disambiguations
	s = strings.Replace(s, " (software)", "", -1)
	s = strings.Replace(s, " (programming language)", "", -1)
	return s
}

func Titles(pages []string) []string {
	var titles []string
	for _, p := range pages {
		titles = append(titles, Title(p))
	}
	return titles
}
