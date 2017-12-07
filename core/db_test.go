package core

import (
	"testing"
	"time"
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

// InterfaceTestDB is used to test DB implementations
func InterfaceTestDB(t *testing.T, db DB) {
	var (
		now     = time.Now().UTC().Round(time.Second) // PG timestamps are not very precise
		test1   = "test_1"
		test1_1 = Page{
			Page:          test1,
			T:             now.Add(-time.Hour),
			StableVersion: "1.0",
			Homepage:      "http://test1.example.com",
		}
		test1_2 = Page{
			Page:          test1,
			T:             now.Add(-time.Minute),
			StableVersion: "2.0",
			Homepage:      "https://test1.example.com",
		}
		test1_3 = Page{
			Page:          test1,
			T:             now,
			StableVersion: "2.0",
			Homepage:      "https://test1.example.com",
		}
		test2   = "test_2"
		test2_1 = Page{
			Page:          test2,
			T:             now,
			StableVersion: "1.0",
		}
	)
	for _, p := range []Page{
		test1_1,
		test1_2,
		test1_3,
		test2_1,
	} {
		if err := db.Store(p); err != nil {
			t.Fatal(err)
		}
	}

	{
		l, err := db.Last(test1)
		if err != nil {
			t.Fatal(err)
		}
		if l == nil {
			t.Fatal("unexpected nil")
		}
		if have, want := *l, test1_3; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
	}

	{
		ls, err := db.Current(test1)
		if err != nil {
			t.Fatal(err)
		}
		if have, want := len(ls), 1; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
		// test1_2 is the most recent spider with a change
		if have, want := ls[0], test1_2; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
	}
}
