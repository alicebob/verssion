package core

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestWikiBasePage(t *testing.T) {
	for url, page := range map[string]string{
		"https://en.wikipedia.org/wiki/Npm":                        "Npm",
		"https://en.wikipedia.org/wiki/Npm_(software)":             "Npm_(software)",
		"https://en.wikipedia.org/wiki/Npm_(software)#Description": "Npm_(software)",
		"/wiki/Npm":                          "Npm",
		"Npm":                                "Npm",
		"Npm (software)":                     "Npm_(software)",
		"/wiki/Npm (software)":               "Npm_(software)",
		"https://en.wikipedia.org/wiki/A/UX": "A/UX",
	} {
		if have, want := WikiBasePage(url), page; have != want {
			t.Errorf("have %q, want %q", have, want)
		}
	}
}

func TestWikipediaSpider(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	w, cb := FixedWikiServer(map[string]string{
		"blue": "a color",
		"red":  "also a color",
		"debian": `<table>
					<tr><td>Latest release</td><td>Version 1</td></tr>
					<tr><td>Website</td><td><a href="http://debian.org">debian.org</a></td></tr>
		</table>`,
		"leftpad": mustRead(t, "./data/leftpad.html"),
	})
	defer w.Close()
	s := NewWikipediaSpider(cb)

	{
		_, err := s.Spider("blue")
		if err == nil {
			t.Fatal("expected an error")
		}
		if _, ok := err.(ErrNoVersion); !ok {
			t.Fatalf("have %T, want ErrNoVersion", err)
		}
	}

	{
		_, err := s.Spider("no such")
		if err == nil {
			t.Fatal("expected an error")
		}
		if _, ok := err.(ErrNotFound); !ok {
			t.Fatalf("have %T, want ErrNotFound", err)
		}
	}

	{
		debian, err := s.Spider("debian")
		if err != nil {
			t.Fatal(err)
		}
		if have, want := debian.StableVersion, "Version 1"; have != want {
			t.Errorf("have %q, want %q", have, want)
		}
		if have, want := debian.Homepage, "[debian.org](http://debian.org)"; have != want {
			t.Errorf("have %q, want %q", have, want)
		}
	}

	{
		_, err := s.Spider("leftpad")
		if err == nil {
			t.Fatal("expected an error")
		}
		redir, ok := err.(ErrRedirect)
		if !ok {
			t.Fatalf("have %T, want ErrNotFound", err)
		}
		if have, want := redir.Page, "leftpad"; have != want {
			t.Errorf("have %q, want %q", have, want)
		}
		if have, want := redir.To, "Npm_(software)"; have != want {
			t.Errorf("have %q, want %q", have, want)
		}
	}
}

// FixedWikiServer is a test helper to have something to spider. It returns a
// server and a callback for core.NewSpider()
func FixedWikiServer(pages map[string]string) (*httptest.Server, func(string) string) {
	r := httprouter.New()
	r.GET("/pages/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		page := p.ByName("id")
		body, ok := pages[page]
		if !ok {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		w.Write([]byte(body))
	})
	s := httptest.NewServer(r)
	cb := func(page string) string {
		return s.URL + "/pages/" + page
	}
	return s, cb
}

func mustRead(t *testing.T, p string) string {
	c, err := ioutil.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return string(c)
}
