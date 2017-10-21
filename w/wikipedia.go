package w

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Revision struct {
	PageID int
	RevID  int
	T      time.Time
}

// Get the latest revisions for the given pages. Calls wikipedia
func Revisions(pages []string) (map[string]Revision, error) {
	p, err := http.Get(wikiURL(
		"format", "json",
		"action", "query",
		"titles", strings.Join(pages, "|"),
		"prop", "revisions",
		"rvprops", "ids|timestamp",
	))
	if err != nil {
		return nil, err
	}
	defer p.Body.Close()

	var pg struct {
		Query struct {
			Pages map[string]struct {
				PageID    int    `json:"pageid"`
				Title     string `json:"title"`
				Revisions []struct {
					RevID    int       `json:"revid"`
					ParentID int       `json:"parentid"`
					T        time.Time `json:"timestamp"`
				} `json:"revisions"`
			} `json:"pages"`
		} `json:"query"`
	}
	if err := json.NewDecoder(p.Body).Decode(&pg); err != nil {
		return nil, err
	}

	res := map[string]Revision{}
	for _, p := range pg.Query.Pages {
		if len(p.Revisions) < 1 {
			continue
		}
		r := p.Revisions[0]
		res[p.Title] = Revision{
			PageID: p.PageID,
			RevID:  r.RevID,
			T:      r.T,
		}
	}
	return res, nil
}

// Get the (XML) parsetree from en.wikipedia
func GetParseTree(rev int) (string, error) {
	p, err := http.Get(wikiURL(
		"format", "json",
		"action", "parse",
		"oldid", strconv.Itoa(rev),
		"prop", "parsetree",
	))
	if err != nil {
		return "", err
	}
	defer p.Body.Close()

	var pt struct {
		Parse struct {
			Title     string            `json:"title"`
			PageID    int               `json:"pageid"`
			ParseTree map[string]string `json:"parsetree"`
		} `json:"parse"`
	}
	err = json.NewDecoder(p.Body).Decode(&pt)
	return pt.Parse.ParseTree["*"], err
}

// args need to be pairs. Or panics.
func wikiURL(args ...string) string {
	v := url.Values{}
	for i := 0; i < len(args); i += 2 {
		v.Add(args[i], args[i+1])
	}
	u := &url.URL{
		Scheme:   "https",
		Host:     "en.wikipedia.org",
		Path:     "/w/api.php",
		RawQuery: v.Encode(),
	}
	return u.String()
}
