package w

import (
	"io"
	"net/http"
	"strings"
)

type Page struct {
	StableVersion string
}

// GetPage downloads and parses given wikipage
func GetPage(title string) (Page, error) {
	p, err := http.Get(wikiURL(title))
	if err != nil {
		return Page{}, err
	}
	defer p.Body.Close()

	return Page{
		StableVersion: StableVersion(p.Body),
	}, nil
}

func StableVersion(n io.Reader) string {
	ts, err := FindTables(n)
	if err != nil {
		return ""
	}
	for _, t := range ts {
		for _, r := range t.Rows {
			if len(r) < 2 {
				continue
			}
			k, v := r[0], r[1]
			switch k {
			case "Stable release", "Latest release":
				return strings.Split(v, ";")[0]
			}
		}
	}
	return ""
}

func wikiURL(page string) string {
	return "https://en.wikipedia.org/wiki/" + page
}
