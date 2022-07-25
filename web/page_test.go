package web_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/verssion/core"
	"github.com/alicebob/verssion/internal"
	"github.com/alicebob/verssion/web"
)

func TestPages(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	var (
		db     = core.NewPGX(internal.TestDB(t))
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
	with(t, body,
		mustcontain("<title>Pages overview"),
		mustcontain("Debian"),
		mustcontain("Glasgow Haskell Compiler"),
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
		db     = core.NewPGX(internal.TestDB(t))
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
	with(t, body,
		mustcontain("<title>Foo"),
		mustcontain("Version 1"),
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
		location := getRedirect(t, s, "/p/Glasgow_Haskell_Compiler")
		if have, want := location, "/p/Glasgow_Haskell_Compiler/"; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
	}

	status, body := get(t, s, "/p/Glasgow_Haskell_Compiler/")
	if have, want := status, 200; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}
	with(t, body,
		mustcontain("<title>Glasgow Haskell Compiler"),
		mustcontain("Glasgow Haskell Compiler"),
		mustcontain("https://haskell.org"),
		mustcontain("8.2.1"),
		mustcontain("8.2.0"),
	)

	{
		status, body := get(t, s, "/p/nosuch/")
		if have, want := status, 404; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
		with(t, body,
			mustcontain("not found"),
		)
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
		with(t, body,
			mustcontain("<title>OS/2"),
			mustcontain("4.52"),
		)
	}

	{
		status, body := get(t, s, "/p/Android/")
		if have, want := status, 200; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
		with(t, body,
			mustcontain("<title>Android"),
			mustcontain("No version"),
		)
	}

	{
		status, body := get(t, s, "/p/Glasgow_Haskell_Compiler/?format=json")
		if have, want := status, 200; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
		var res web.PageJSON
		if err := json.Unmarshal([]byte(body), &res); err != nil {
			t.Fatal(err)
		}
		if have, want := res.Page, "Glasgow_Haskell_Compiler"; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
		if have, want := res.Title, "Glasgow Haskell Compiler"; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
		if have, want := res.StableVersion, "8.2.1 / July 22, 2017"; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
		if have, want := res.Homepage, "https://haskell.org/ghc"; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
	}

	{
		status, _ := get(t, s, "/p/Glasgow_Haskell_Compiler/?format=foo")
		if have, want := status, 400; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
	}

	{
		status, _ := get(t, s, "/p/Android/?format=json")
		if have, want := status, 404; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
	}

}
