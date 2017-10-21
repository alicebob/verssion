package w

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// Get the (XML) parsetree from en.wikipedia
func GetParseTree(page string) (string, error) {
	p, err := http.Get(wikiURL(page))
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

func wikiURL(page string) string {
	v := url.Values{}
	v.Add("format", "json")
	v.Add("page", page)
	v.Add("action", "parse")
	v.Add("prop", "parsetree")
	u := &url.URL{
		Scheme:   "https",
		Host:     "en.wikipedia.org",
		Path:     "/w/api.php",
		RawQuery: v.Encode(),
	}
	return u.String()
}
