package core

import (
	"reflect"
	"testing"
	"time"
)

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
			T:             now.Add(time.Second),
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
		ls, err := db.Current(0, SpiderT)
		if err != nil {
			t.Fatal(err)
		}
		if have, want := len(ls), 2; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
		if have, want := ls[0], test2_1; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
		if have, want := ls[1], test1_2; have != want {
			t.Fatalf("have %v, want %v", have, want)
		}
	}

	{
		ls, err := db.CurrentIn(test1)
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

	{
		l, err := db.Last("nosuchpage")
		if err != nil {
			t.Fatal(err)
		}
		if l != nil {
			t.Fatalf("not a nil: %v", l)
		}
	}
}

// InterfaceTestCurated is used to test the Curated methods of DB implementations
func InterfaceTestCurated(t *testing.T, db DB) {
	{
		c, err := db.LoadCurated("nosuch")
		if err != nil {
			t.Fatal(err)
		}
		if c != nil {
			t.Fatal("want nil")
		}
	}
	id, err := db.CreateCurated()
	if err != nil {
		t.Fatal(err)
	}
	if have, want := len(id), 36; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}
	if db.CuratedSetPages(id, []string{"page1", "page2"}); err != nil {
		t.Fatal(err)
	}
	if db.CuratedSetPages(id, []string{"page3", "page2"}); err != nil {
		t.Fatal(err)
	}
	if db.CuratedSetTitle(id, "My first list"); err != nil {
		t.Fatal(err)
	}

	c, err := db.LoadCurated(id)
	if err != nil {
		t.Fatal(err)
	}
	if have, want := c.Pages, []string{"page2", "page3"}; !reflect.DeepEqual(have, want) {
		t.Fatalf("have %#v, want %#v", have, want)
	}
	if have, want := c.CustomTitle, "My first list"; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}
	if have, want := c.Title(), "My first list"; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}
	if c.Created.IsZero() {
		t.Fatal("0 created")
	}
	if c.LastUpdated.IsZero() {
		t.Fatal("0 last updated")
	}
	if have, want := c.LastUsed, c.Created; !want.Equal(have) {
		t.Fatalf("have %v, want %v", have, want)
	}

	// Check SetUsed
	{
		if db.CuratedSetUsed(id); err != nil {
			t.Fatal(err)
		}
		c, err := db.LoadCurated(id)
		if err != nil {
			t.Fatal(err)
		}
		if have, want := c.LastUsed, c.Created; want.Equal(have) {
			t.Fatalf("not: have %v, want %v", have, want)
		}
	}

	{
		id2, err := db.CreateCurated()
		if err != nil {
			t.Fatal(err)
		}
		if id == id2 {
			t.Fatalf("double curated ID")
		}
	}

	if have, want := db.CuratedSetPages("nosuch", nil), ErrCuratedNotFound; have != want {
		t.Fatalf("have error %v, want error %v", have, want)
	}
	if have, want := db.CuratedSetTitle("nosuch", "foo"), ErrCuratedNotFound; have != want {
		t.Fatalf("have error %v, want error %v", have, want)
	}
	if have, want := db.CuratedSetUsed("nosuch"), ErrCuratedNotFound; have != want {
		t.Fatalf("have error %v, want error %v", have, want)
	}
}
