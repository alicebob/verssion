package web_test

import (
	"io/ioutil"
	"log"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/verssion/core"
	"github.com/alicebob/verssion/web"
)

func TestPages(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	var (
		db     = core.NewMemory()
		spider = NewFixedSpider()
		m      = web.Mux("", db, spider, "")
	)
	db.Store(core.Page{
		Page:          "Debian",
		StableVersion: "my version",
		T:             time.Now(),
	})
	db.Store(core.Page{
		Page:          "Glasgow_Haskell_Compiler",
		StableVersion: "8.2.1 / July 22, 2017",
		T:             time.Now().Add(time.Minute),
	})
	s := httptest.NewServer(m)
	defer s.Close()

	{
		location := getRedirect(t, s, "/p")
		if have, want := location, "/p/"; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
	}

	status, body := get(t, s, "/p/")
	if have, want := status, 200; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}
	contains(t, body,
		"<title>Pages overview",
		"Debian",
		"Glasgow Haskell Compiler",
	)

	{
		status, sbody := get(t, s, "/p/?order=spider")
		if have, want := status, 200; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
		if sbody == body {
			t.Fatal("order change doesn't do anything")
		}
	}
}

func TestNewPage(t *testing.T) {
	var (
		db     = core.NewMemory()
		spider = NewFixedSpider(core.Page{
			Page:          "Foo",
			StableVersion: "Version 1",
			T:             time.Now(),
		})
		m = web.Mux("", db, spider, "")
	)
	s := httptest.NewServer(m)
	defer s.Close()

	{
		location := getRedirect(t, s, "/p/?page=https://en.wikipedia.org/wiki/Foo")
		if have, want := location, "/p/Foo/"; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
	}

	status, body := get(t, s, "/p/Foo/")
	if have, want := status, 200; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}
	contains(t, body,
		"<title>Foo",
		"Version 1",
	)
}

func TestPage(t *testing.T) {
	var (
		db     = core.NewMemory()
		spider = NewFixedSpider()
		m      = web.Mux("", db, spider, "")
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
		StableVersion: "8.2.0",
		T:             time.Now(),
	})
	db.Store(core.Page{
		Page:          "Glasgow_Haskell_Compiler",
		StableVersion: "8.2.1 / July 22, 2017",
		Homepage:      "https://haskell.org/ghc",
		T:             time.Now(),
	})
	db.Store(core.Page{
		Page:          "OS/2",
		StableVersion: "4.52 (2001)",
		T:             time.Now(),
	})
	spider.SetError("Android", core.ErrNoVersion{Page: "Android"})

	{
		status, _ := get(t, s, "/p/Glasgow_Haskell_Compiler")
		if have, want := status, 301; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
	}

	status, body := get(t, s, "/p/Glasgow_Haskell_Compiler/")
	if have, want := status, 200; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}
	contains(t, body,
		"<title>Glasgow Haskell Compiler",
		"Glasgow Haskell Compiler",
		"https://haskell.org",
		"8.2.1",
		"8.2.0",
	)

	{
		status, body := get(t, s, "/p/nosuch/")
		if have, want := status, 404; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
		contains(t, body, "not found")
	}

	{
		status, _ := get(t, s, "/p/OS/2")
		if have, want := status, 301; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
	}

	{
		status, body := get(t, s, "/p/OS/2/")
		if have, want := status, 200; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
		contains(t, body,
			"<title>OS/2",
			"4.52",
		)
	}

	{
		status, body := get(t, s, "/p/Android/")
		if have, want := status, 200; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
		contains(t, body,
			"<title>Android",
			"No version",
		)
	}
}
