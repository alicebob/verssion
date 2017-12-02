package w

import (
	"testing"
)

func TestDefaultTitle(t *testing.T) {
	type cas struct {
		Pages []string
		Want  string
	}
	for i, c := range []cas{
		{
			Want: "[untitled feed]",
		},
		{
			Pages: []string{"Foo"},
			Want:  "Foo",
		},
		{
			Pages: []string{"Foo", "Bar", "Baz", "Baq"},
			Want:  "Foo, Bar, Baz, Baq",
		},
		{
			Pages: []string{"Foo", "Bar", "Baz", "and", "even", "more"},
			Want:  "Foo, Bar, Baz, and, ... (2 more)",
		},
	} {
		cur := Curated{
			Pages: c.Pages,
		}
		if have, want := cur.DefaultTitle(), c.Want; have != want {
			t.Errorf("case %d: have: %v, want: %v", i, have, want)
		}
	}
}
