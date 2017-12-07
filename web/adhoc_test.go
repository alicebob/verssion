package web_test

import (
	"encoding/xml"
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
		m  = web.Mux("", db, web.NotFetcher(), "")
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
