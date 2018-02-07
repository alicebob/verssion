package web_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alicebob/verssion/core"
	"github.com/alicebob/verssion/web"
)

func TestIndex(t *testing.T) {
	var (
		db = core.NewMemory()
		m  = web.Mux("/", db, nil, "")
	)
	s := httptest.NewServer(m)
	defer s.Close()
	db.Store(core.Page{Page: "Debian", StableVersion: "my version"})
	db.Store(core.Page{Page: "Glasgow_Haskell_Compiler", StableVersion: "8.2.1 / July 22, 2017"})
	db.Store(core.Page{Page: "WithLink", StableVersion: "a [link](https://link.me)!"})

	status, body := get(t, s, "")
	if have, want := status, 200; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	with(t, body,
		mustcontain("<title>VeRSSion"),
		mustcontain("Debian"),
		mustcontain("my version"),
		mustcontain("Glasgow Haskell Compiler"), // titlefication test
		mustcontain("Glasgow_Haskell_Compiler"),
		mustcontain("a link!"),
	)
}
