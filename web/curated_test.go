package web_test

import (
	"encoding/xml"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"
	"time"

	"github.com/alicebob/verssion/core"
	"github.com/alicebob/verssion/web"
)

func TestCurated(t *testing.T) {
	var (
		db = core.NewMemory()
		m  = web.Mux("/", db, NewFixedSpider(), "")
	)
	s := httptest.NewServer(m)
	defer s.Close()
	db.Store(core.Page{
		Page:          "Debian",
		StableVersion: "my version",
		T:             time.Now(),
	})
	db.Store(core.Page{
		Page:          "Glasgow_Haskell_Compiler",
		StableVersion: "8.2.1 / July 22, 2017",
		T:             time.Now(),
	})

	{
		status, body := get(t, s, "/curated/")
		if have, want := status, 200; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
		contains(t, body,
			"<title>New feed",
			"<h2>New feed",
			"Glasgow Haskell Compiler",
			"Glasgow_Haskell_Compiler",
		)
	}

	var curURL string
	{
		c := s.Client()
		c.CheckRedirect = noRedirect
		base := "/curated/"
		r, err := c.PostForm(s.URL+base, url.Values{
			"p": []string{"Debian"},
		})
		if err != nil {
			t.Fatal(err)
		}
		r.Body.Close()
		if have, want := r.StatusCode, 302; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
		redir := r.Header.Get("Location")
		curURL = path.Join(base, redir) + "/"
	}

	{
		status, body := get(t, s, curURL)
		if have, want := status, 200; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
		contains(t, body,
			"Debian",
		)
	}

	// Should be a cookie with the feed on the index page
	{
		status, body := get(t, s, "")
		if have, want := status, 200; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
		contains(t, body,
			"Your recent feeds",
			"Debian",
		)
	}

	{
		status, body := get(t, s, curURL+"atom.xml")
		if have, want := status, 200; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
		contains(t, body,
			"Debian",
		)

		var f web.Feed
		if err := xml.Unmarshal([]byte(body), &f); err != nil {
			t.Fatal(err)
		}
		if have, want := len(f.Entries), 1; have != want {
			t.Errorf("have %v, want %v", have, want)
		}
	}

	{
		status, body := get(t, s, curURL+"edit.html")
		if have, want := status, 200; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
		contains(t, body,
			"Debian",
		)
	}
}
