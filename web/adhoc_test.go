package web_test

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/verssion/core"
	"github.com/alicebob/verssion/internal"
	"github.com/alicebob/verssion/web"
)

func TestAdhoc(t *testing.T) {
	core.CacheOK = 42 * 365 * 24 * time.Hour // a while
	var (
		db  = core.NewPGX(internal.TestDB(t))
		m   = web.Mux("", db, NewFixedSpider(), "")
		now = time.Date(2017, 3, 14, 15, 14, 0, 0, time.UTC)
	)
	s := httptest.NewServer(m)
	defer s.Close()
	db.Store(core.Page{Page: "Debian", StableVersion: "my version"})
	db.Store(core.Page{Page: "Glasgow_Haskell_Compiler", StableVersion: "8.1.0 / July 20, 2015", T: now})
	db.Store(core.Page{Page: "Glasgow_Haskell_Compiler", StableVersion: "8.2.0 / July 21, 2016", T: now.Add(time.Second)})
	db.Store(core.Page{Page: "Glasgow_Haskell_Compiler", StableVersion: "8.2.1 / July 22, 2017", T: now.Add(2 * time.Second)})

	status, body := get(t, s, "/adhoc/atom.xml?p=Glasgow_Haskell_Compiler")
	if have, want := status, 200; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}
	with(t, body,
		mustcontain("Glasgow Haskell Compiler"),
		mustcontain("<content>8.2.1"),
	)

	var f web.Feed
	if err := xml.Unmarshal([]byte(body), &f); err != nil {
		t.Fatal(err)
	}
	if have, want := len(f.Entries), 3; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestAdhoc404(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	var (
		db = core.NewPGX(internal.TestDB(t))
		m  = web.Mux("", db, NewFixedSpider(), "")
	)
	s := httptest.NewServer(m)
	defer s.Close()

	status, _ := get(t, s, "/adhoc/atom.xml?p=Foobar")
	if have, want := status, 404; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}
}

func TestAdhocNoLink(t *testing.T) {
	var (
		db = core.NewPGX(internal.TestDB(t))
		m  = web.Mux("", db, NewFixedSpider(), "")
	)
	s := httptest.NewServer(m)
	defer s.Close()
	db.Store(core.Page{
		Page:          "Z_shell",
		StableVersion: "[5.4.2](https://sourceforge.net/projects/zsh/files/zsh/5.4.2/) / August 28, 2017",
		T:             time.Now(),
	})

	status, body := get(t, s, "/adhoc/atom.xml?p=Z_shell")
	if have, want := status, 200; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}
	with(t, body,
		mustcontain("Z shell"),
		mustcontain("5.4.2"),
		mustnotcontain("sourceforge"),
	)
}
