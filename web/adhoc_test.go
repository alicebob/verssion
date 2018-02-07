package web_test

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/verssion/core"
	"github.com/alicebob/verssion/web"
)

func TestAdhoc(t *testing.T) {
	var (
		db = core.NewMemory()
		m  = web.Mux("", db, NewFixedSpider(), "")
	)
	s := httptest.NewServer(m)
	defer s.Close()
	db.Store(core.Page{Page: "Debian", StableVersion: "my version"})
	db.Store(core.Page{Page: "Glasgow_Haskell_Compiler", StableVersion: "8.1.0 / July 20, 2015", T: time.Now()})
	db.Store(core.Page{Page: "Glasgow_Haskell_Compiler", StableVersion: "8.2.0 / July 21, 2016", T: time.Now()})
	db.Store(core.Page{Page: "Glasgow_Haskell_Compiler", StableVersion: "8.2.1 / July 22, 2017", T: time.Now()})

	status, body := get(t, s, "/adhoc/atom.xml?p=Glasgow_Haskell_Compiler")
	if have, want := status, 200; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}
	if in, want := body, "Glasgow Haskell Compiler"; !strings.Contains(in, want) {
		t.Fatalf("no %q found", want)
	}
	if in, want := body, "<content>8.2.1 "; !strings.Contains(in, want) {
		t.Fatalf("no %q found", want)
	}

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
		db = core.NewMemory()
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
		db = core.NewMemory()
		m  = web.Mux("", db, NewFixedSpider(), "")
	)
	s := httptest.NewServer(m)
	defer s.Close()
	db.Store(core.Page{Page: "Z_shell", StableVersion: "[5.4.2](https://sourceforge.net/projects/zsh/files/zsh/5.4.2/) / August 28, 2017", T: time.Now()})

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
