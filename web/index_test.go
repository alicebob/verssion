package web_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alicebob/verssion/core"
	"github.com/alicebob/verssion/web"
)

func TestIndex(t *testing.T) {
	var (
		db = core.NewMemory()
		m  = web.Mux("/", db, web.NotFetcher(), "")
	)
	s := httptest.NewServer(m)
	defer s.Close()
	db.Store(core.Page{Page: "Debian", StableVersion: "my version"})
	db.Store(core.Page{Page: "Glasgow_Haskell_Compiler", StableVersion: "8.2.1 / July 22, 2017"})

	status, body := get(t, s, "")
	if have, want := status, 200; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	if in, want := body, "<title>VeRSSion"; !strings.Contains(in, want) {
		t.Fatalf("no %q found", want)
	}

	if in, want := body, "Debian"; !strings.Contains(in, want) {
		t.Fatalf("no %q found", want)
	}
	if in, want := body, "my version"; !strings.Contains(in, want) {
		t.Fatalf("no %q found", want)
	}
	// titleification test
	if in, want := body, "Glasgow Haskell Compiler"; !strings.Contains(in, want) {
		t.Fatalf("no %q found", want)
	}
	if in, want := body, "Glasgow_Haskell_Compiler/"; !strings.Contains(in, want) {
		t.Fatalf("no %q found", want)
	}
}
